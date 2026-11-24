package service

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    jwtauth "github.com/herenote/backend/internal/auth"
    "github.com/herenote/backend/internal/model"
    "github.com/herenote/backend/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

// 위치 인증
type AuthService struct {
    sessionRepo *repository.SessionRepository
    placeRepo *repository.PlaceRepository
}

func NewAuthService(s *repository.SessionRepository, p *repository.PlaceRepository) *AuthService {
    return &AuthService{sessionRepo: s, placeRepo: p}
}

func (s *AuthService) Verify(ctx context.Context, placeID string, req *model.VerifyRequest) (*model.VerifyResponse, error) {
    // PostGIS로 서버 측 위치 검증
    verified, err := s.placeRepo.VerifyLocation(ctx, placeID, req.Latitude, req.Longitude)
    if err != nil {
        return nil, fmt.Errorf("장소 검증 실패: %w", err)
    }

    if !verified {
        return &model.VerifyResponse{
            Verified: false,
            Message: "해당 장소에 도달하지 못했습니다. 현장에서 다시 시도해주세요.",
        }, nil
    }

    // 인증 성공 - 24시간 유효 토큰 발급
    token := uuid.NewString()
    expiredAt := time.Now().Add(24 * time.Hour)

    if err := s.sessionRepo.Create(ctx, token, placeID, req.UserID, expiredAt); err != nil {
        return nil, fmt.Errorf("세션 생성 실패:%w", err)
    }

    return &model.VerifyResponse{
        Verified: true,
        Token: token,
        Message: "위치 인증 성공, 방명록을 열람하고 작성할 수 있습니다.",
    }, nil
}

// 장소

type PlaceService struct {
    repo *repository.PlaceRepository
}

func NewPlaceService(r *repository.PlaceRepository) *PlaceService {
    return &PlaceService{repo: r}
}

func (s *PlaceService) List(ctx context.Context) ([]model.Place, error) {
    return s.repo.FindAll(ctx)
}

func (s *PlaceService) Get(ctx context.Context, id string) (*model.Place, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *PlaceService) Create(ctx context.Context, req *model.CreatePlaceRequest) (*model.Place, error) {
    return s.repo.Create(ctx, req)
}

func (s *PlaceService) FindNearby(ctx context.Context, lat, lng float64) ([]model.Place, error) {
    return s.repo.FindNearby(ctx, lat, lng)
}

// 방명록
// userFinder는 작성자의 현재 뱃지를 DB에서 읽기 위한 최소 인터페이스입니다.
type userFinder interface {
    FindByID(ctx context.Context, id string) (*model.User, error)
}

type GuestbookService struct {
    repo        *repository.GuestbookRepository
    sessionRepo *repository.SessionRepository
    users       userFinder
}

func NewGuestbookService(g *repository.GuestbookRepository, s *repository.SessionRepository, u userFinder) *GuestbookService {
    return &GuestbookService{repo: g, sessionRepo: s, users: u}
}

func (s *GuestbookService) List(ctx context.Context, placeID string) ([]model.Guestbook, error) {
    return s.repo.ListByPlace(ctx, placeID)
}

func (s *GuestbookService) Create(ctx context.Context, placeID, userID string, req *model.CreateGuestbookRequest) (*model.Guestbook, error) {
    // 뱃지 보유자는 자동 상단 고정.
    // 예전에는 위치 인증 미들웨어가 넘겨준 badge_type을 썼는데, 미들웨어가 항상
    // NONE을 넣고 있어서 고정이 한 번도 동작하지 않았다. 작성 시점에 DB에서 직접 읽는다.
    return s.repo.Create(ctx, placeID, userID, req, s.shouldPin(ctx, userID))
}

// shouldPin은 작성자가 뱃지를 갖고 있는지 확인합니다.
// 조회에 실패하면 글 작성은 막지 않고 고정만 하지 않습니다.
func (s *GuestbookService) shouldPin(ctx context.Context, userID string) bool {
    if s.users == nil || userID == "" {
        return false
    }
    u, err := s.users.FindByID(ctx, userID)
    if err != nil || u == nil {
        return false
    }
    return u.BadgeType != "" && u.BadgeType != model.BadgeNone
}

func (s *GuestbookService) MyPlaces(ctx context.Context, userID string) ([]model.Place, error) {
    return s.repo.FindMyPlaces(ctx, userID)
}

func (s *GuestbookService) UserGuestbooks(ctx context.Context, userID, placeID string) ([]model.Guestbook, error) {
    return s.repo.FindByUserAndPlace(ctx, userID, placeID)
}

func (s *GuestbookService) Stats(ctx context.Context, userID string) (visited, notes int, err error) {
    places, e := s.repo.FindMyPlaces(ctx, userID)
    if e != nil {
        return 0, 0, e
    }
    notes, e = s.repo.CountByUser(ctx, userID)
    return len(places), notes, e
}

// 유저
type UserService struct {
    repo      *repository.UserRepository
    jwtSecret string
}

func NewUserService(r *repository.UserRepository, jwtSecret string) *UserService {
    return &UserService{repo: r, jwtSecret: jwtSecret}
}

func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.AuthResponse, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user, err := s.repo.Create(ctx, req, string(hash))
    if err != nil {
        return nil, err
    }

    token, err := jwtauth.Sign(user.ID, user.Nickname, s.jwtSecret)
    if err != nil {
        return nil, fmt.Errorf("토큰 발급 실패: %w", err)
    }

    return &model.AuthResponse{
        ID:        user.ID,
        Nickname:  user.Nickname,
        BadgeType: user.BadgeType,
        Token:     token,
    }, nil
}

func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
    user, err := s.repo.FindByNickname(ctx, req.Nickname)
    if err != nil {
        return nil, fmt.Errorf("닉네임 또는 비밀번호가 올바르지 않습니다.")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return nil, fmt.Errorf("닉네임 또는 비밀번호가 올바르지 않습니다.")
    }

    token, err := jwtauth.Sign(user.ID, user.Nickname, s.jwtSecret)
    if err != nil {
        return nil, fmt.Errorf("토큰 발급 실패: %w", err)
    }

    return &model.AuthResponse{
        ID:        user.ID,
        Nickname:  user.Nickname,
        BadgeType: user.BadgeType,
        Token:     token,
    }, nil
}

func (s *UserService) Get(ctx context.Context, id string) (*model.User, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *UserService) GrantBadge(ctx context.Context, userID string, badge model.BadgeType) error {
    return s.repo.GrantBadge(ctx, userID, badge)
}

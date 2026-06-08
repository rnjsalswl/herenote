package service

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/herenote/backend/internal/model"
    "github.com/herenote/backend/internal/repository"
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

// 방명록
type GuestbookService struct {
    repo *repository.GuestbookRepository
    sessionRepo *repository.SessionRepository
}

func NewGuestbookService(g *repository.GuestbookRepository, s *repository.SessionRepository) *GuestbookService {
    return &GuestbookService{repo: g, sessionRepo: s}
}

func (s *GuestbookService) List(ctx context.Context, placeID string) ([]model.Guestbook, error) {
    return s.repo.ListByPlace(ctx, placeID)
}

func (s *GuestbookService) Create(ctx context.Context, placeID, userID string, req *model.CreateGuestbookRequest, badgeType model.BadgeType) (*model.Guestbook, error) {
    // 뱃지 보유자는 자동 상단 고정
    isPinned := badgeType != model.BadgeNone
    return s.repo.Create(ctx, placeID, userID, req, isPinned)
}

// 유저
type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
    return &UserService{repo: r}
}

func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
    return s.repo.Create(ctx, req)
}

func (s *UserService) Get(ctx context.Context, id string) (*model.User, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *UserService) GrantBadge(ctx context.Context, userID string, badge model.BadgeType) error {
    return s.repo.GrantBadge(ctx, userID, badge)
}
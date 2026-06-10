package repository

import (
    "context"

    "github.com/herenote/backend/internal/model"
    "github.com/jackc/pgx/v5/pgxpool"
)

type GuestbookRepository struct {
    db *pgxpool.Pool
}

func NewGuestbookRepository(db *pgxpool.Pool) *GuestbookRepository {
    return &GuestbookRepository{db: db}
}

// 뱃지 보유자 상단 고정 후 최신순
func (r *GuestbookRepository) ListByPlace(ctx context.Context, placeID string) ([]model.Guestbook, error) {
    rows, err := r.db.Query(ctx, `
        SELECT g.id, g.place_id, g.user_id,
               u.nickname, u.badge_type,
               g.content, COALESCE(g.image_url, ''),
               g.is_pinned, g.created_at
        FROM guestbooks g
        JOIN users u ON g.user_id = u.id
        WHERE g.place_id = $1
        ORDER BY g.is_pinned DESC, g.created_at DESC
    `, placeID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.Guestbook
    for rows.Next() {
        var g model.Guestbook
        if err := rows.Scan(
            &g.ID, &g.PlaceID, &g.UserID,
            &g.Nickname, &g.BadgeType,
            &g.Content, &g.ImageURL,
            &g.IsPinned, &g.CreatedAt,
        ); err != nil {
            return nil, err
        }
        list = append(list, g)
    }
    return list, nil
}

func (r *GuestbookRepository) Create(ctx context.Context, placeID, userID string, req *model.CreateGuestbookRequest, isPinned bool) (*model.Guestbook, error) {
    var g model.Guestbook
    err := r.db.QueryRow(ctx, `
        INSERT INTO guestbooks (place_id, user_id, content, image_url, is_pinned)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, place_id, user_id,
                  content, COALESCE(image_url, ''),
                  is_pinned, created_at
    `, placeID, userID, req.Content, req.ImageURL, isPinned,
    ).Scan(
        &g.ID, &g.PlaceID, &g.UserID,
        &g.Content, &g.ImageURL,
        &g.IsPinned, &g.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &g, nil
}

// 특정 장소에서 특정 유저가 쓴 방명록
func (r *GuestbookRepository) FindByUserAndPlace(ctx context.Context, userID, placeID string) ([]model.Guestbook, error) {
    rows, err := r.db.Query(ctx, `
        SELECT g.id, g.place_id, g.user_id,
               u.nickname, u.badge_type,
               g.content, COALESCE(g.image_url, ''),
               g.is_pinned, g.created_at
        FROM guestbooks g
        JOIN users u ON g.user_id = u.id
        WHERE g.place_id = $1 AND g.user_id = $2
        ORDER BY g.created_at DESC
    `, placeID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.Guestbook
    for rows.Next() {
        var g model.Guestbook
        if err := rows.Scan(
            &g.ID, &g.PlaceID, &g.UserID,
            &g.Nickname, &g.BadgeType,
            &g.Content, &g.ImageURL,
            &g.IsPinned, &g.CreatedAt,
        ); err != nil {
            return nil, err
        }
        list = append(list, g)
    }
    return list, nil
}

// 유저가 작성한 방명록 수
func (r *GuestbookRepository) CountByUser(ctx context.Context, userID string) (int, error) {
    var count int
    err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM guestbooks WHERE user_id = $1`, userID).Scan(&count)
    return count, err
}

// 내가 작성한 방명록의 장소 목록 (최근순, 장소 중복 제거)
func (r *GuestbookRepository) FindMyPlaces(ctx context.Context, userID string) ([]model.Place, error) {
    rows, err := r.db.Query(ctx, `
        SELECT DISTINCT ON (p.id)
               p.id, p.name, p.description,
               ST_Y(p.location::geometry) AS latitude,
               ST_X(p.location::geometry) AS longitude,
               p.radius_meters, p.created_at
        FROM guestbooks g
        JOIN places p ON g.place_id = p.id
        WHERE g.user_id = $1
        ORDER BY p.id, g.created_at DESC
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var places []model.Place
    for rows.Next() {
        var p model.Place
        if err := rows.Scan(
            &p.ID, &p.Name, &p.Description,
            &p.Latitude, &p.Longitude,
            &p.RadiusMeters, &p.CreatedAt,
        ); err != nil {
            return nil, err
        }
        places = append(places, p)
    }
    return places, nil
}
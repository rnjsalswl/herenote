package repository

import (
    "context"
    "time"

    "github.com/herenote/backend/internal/model"
    "github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
    db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
    return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, token, placeID, userID string, expiredAt time.Time) error {
    _, err := r.db.Exec(ctx, `
        INSERT INTO verified_sessions (token, place_id, user_id, expired_at)
        VALUES ($1, $2, $3, $4)
    `, token, placeID, userID, expiredAt)
    return err
}

func (r *SessionRepository) FindValid(ctx context.Context, token, placeID string) (*model.Session, error) {
    var s model.Session
    err := r.db.QueryRow(ctx, `
        SELECT token, place_id, user_id, expired_at
        FROM verified_sessions
        WHERE token = $1
          AND place_id = $2
          AND expired_at > NOW()
    `, token, placeID).Scan(&s.Token, &s.PlaceID, &s.UserID, &s.ExpiredAt)
    if err != nil {
        return nil, err
    }
    return &s, nil
}
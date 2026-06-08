package repository

import (
    "context"

    "github.com/herenote/backend/internal/model"
    "github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, req *model.CreateUserRequest, passwordHash string) (*model.User, error) {
    var u model.User
    err := r.db.QueryRow(ctx, `
        INSERT INTO users (nickname, password_hash)
        VALUES ($1, $2)
        RETURNING id, nickname, badge_type, created_at
    `, req.Nickname, passwordHash).Scan(&u.ID, &u.Nickname, &u.BadgeType, &u.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *UserRepository) FindByNickname(ctx context.Context, nickname string) (*model.User, error) {
    var u model.User
    var passwordHash string
    err := r.db.QueryRow(ctx, `
        SELECT id, nickname, badge_type, created_at, COALESCE(password_hash, '')
        FROM users WHERE nickname = $1
    `, nickname).Scan(&u.ID, &u.Nickname, &u.BadgeType, &u.CreatedAt, &passwordHash)
    if err != nil {
        return nil, err
    }
    u.PasswordHash = passwordHash
    return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
    var u model.User
    err := r.db.QueryRow(ctx, `
        SELECT id, nickname, badge_type, created_at
        FROM users WHERE id = $1
    `, id).Scan(&u.ID, &u.Nickname, &u.BadgeType, &u.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *UserRepository) GrantBadge(ctx context.Context, userID string, badge model.BadgeType) error {
    _, err := r.db.Exec(ctx, `
        UPDATE users SET badge_type = $1 WHERE id = $2
    `, badge, userID)
    return err
}
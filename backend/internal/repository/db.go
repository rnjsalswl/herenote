package repository

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(databaseURL string) (*pgxpool.Pool, error) {
    if databaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL이 설정되지 않음")
    }

    pool, err := pgxpool.New(context.Background(), databaseURL)
    if err != nil {
        return nil, fmt.Errorf("DB 연결 실패: %w", err)
    }

    if err := pool.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("DB ping 실패: %w", err)
    }

    return pool, nil
}
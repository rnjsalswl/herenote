package repository

import (
    "context"

    "github.com/herenote/backend/internal/model"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PlaceRepository struct {
    db *pgxpool.Pool
}

func NewPlaceRepository(db *pgxpool.Pool) *PlaceRepository {
    return &PlaceRepository{db: db}
}

// 전체 장소 목록
func (r *PlaceRepository) FindAll(ctx context.Context) ([]model.Place, error) {
    rows, err := r.db.Query(ctx, `
        SELECT id, name, description,
               ST_Y(location::geometry) AS latitude,
               ST_X(location::geometry) AS longitude,
               radius_meters, created_at
        FROM places
        ORDER BY created_at DESC
    `)
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

// 장소 단건 조회
func (r *PlaceRepository) FindByID(ctx context.Context, id string) (*model.Place, error) {
    var p model.Place
    err := r.db.QueryRow(ctx, `
        SELECT id, name, description,
               ST_Y(location::geometry) AS latitude,
               ST_X(location::geometry) AS longitude,
               radius_meters, created_at
        FROM places WHERE id = $1
    `, id).Scan(
    &p.ID, &p.Name, &p.Description,
    &p.Latitude, &p.Longitude,
    &p.RadiusMeters, &p.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &p, nil
}

// 현재 위치 근처 장소 조회(500m 이내)
func (r *PlaceRepository) FindNearby(ctx context.Context, lat, lng float64) ([]model.Place, error) {
    rows, err := r.db.Query(ctx, `
        SELECT id, name, description,
               ST_Y(location::geometry) AS latitude,
               ST_X(location::geometry) AS longitude,
               radius_meters, created_at
        FROM places
        WHERE ST_DWithin(
            location::geography,
            ST_MakePoint($2, $1)::geography,
            500
        )
        ORDER BY ST_Distance(
            location::geography,
            ST_MakePoint($2, $1)::geography
        )
    `, lat, lng)
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

// GPS 위치 인증 - PostGIS 핵심 쿼리
func (r *PlaceRepository) VerifyLocation(ctx context.Context, placeID string, lat, lng float64) (bool, error) {
    var verified bool
    err := r.db.QueryRow(ctx, `
        SELECT ST_DWithin(
            location::geography,
            ST_MakePoint($2, $3)::geography,
            radius_meters
        )
        FROM places WHERE id = $1
    `, placeID, lng, lat).Scan(&verified)   // ST_MakePoint는 (경도, 위도) 순서
    if err != nil {
        return false, err
    }
    return verified, nil
}

// 장소 생성
func (r *PlaceRepository) Create(ctx context.Context, req *model.CreatePlaceRequest) (*model.Place, error) {
    radius := 100
    if req.RadiusMeters > 0 {
        radius = req.RadiusMeters
    }

    var p model.Place
    err := r.db.QueryRow(ctx, `
        INSERT INTO places (name, description, location, radius_meters)
        VALUES ($1, $2, ST_MakePoint($4, $3)::geography, $5)
        RETURNING id, name, description,
                  ST_Y(location::geometry),
                  ST_X(location::geometry),
                  radius_meters, created_at
    `, req.Name, req.Description, req.Latitude, req.Longitude, radius,
    ).Scan(
        &p.ID, &p.Name, &p.Description,
        &p.Latitude, &p.Longitude,
        &p.RadiusMeters, &p.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &p, nil
}
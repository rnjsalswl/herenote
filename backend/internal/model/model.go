package model

import "time"

// 장소
type Place struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    RadiusMeters int `json:"radius_meters"`
    CreatedAt time.Time `json:"created_at"`
}

type CreatePlaceRequest struct {
    Name string `json:"name"    binding:"required"`
    Description string `json:"description"`
    Latitude float64 `json:"latitude"   binding:"required"`
    Longitude float64 `json:"longitude" binding:"required"`
    RadiusMeters int `json:"radius_meters"`
}

// 유저
type BadgeType string

const (
    BadgeNone BadgeType = "NONE"
    BadgeBlue BadgeType = "BLUE"
    BadgeGold BadgeType = "GOLD"
    BadgeSpecial BadgeType = "SPECIAL"
)

type User struct {
    ID string `json:"id"`
    Nickname string `json:"nickname"`
    BadgeType BadgeType `json:"badge_type"`
    PasswordHash string `json:"-"`  // json 응답에서 숨김
    CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
    Nickname string `json:"nickname" binding:"required"`
    Password string `json:"password" binding:"required,min=4"`
}

type LoginRequest struct {
    Nickname string `json:"nickname" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    UserID string `json:"user_id"`
    Nickname string `json:"nickname"`
    BadgeType BadgeType `json:"badge_type"`
}

type AuthResponse struct {
    ID string `json:"id"`
    Nickname string `json:"nickname"`
    BadgeType BadgeType `json:"badge_type"`
    Token string `json:"token"` // JWT 교체 예정
}

// 방명록
type Guestbook struct {
    ID string `json:"id"`
    PlaceID string `json:"place_id"`
    UserID string `json:"user_id"`
    Nickname string `json:"nickname"`
    BadgeType BadgeType `json:"badge_type"`
    Content string `json:"content"`
    ImageURL string `json:"image_url,omitempty"`
    IsPinned bool `json:"is_pinned"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateGuestbookRequest struct {
	Content  string `json:"content"  binding:"required,max=500"`
	ImageURL string `json:"image_url"`
}

// 위치 인증
type VerifyRequest struct {
    Latitude float64 `json:"latitude"   binding:"required"`
    Longitude float64 `json:"longitude" binding:"required"`
    UserID string `json:"user_id"   binding:"required"`
}

type VerifyResponse struct {
    Verified bool `json:"verified"`
    Token string `json:"token,omitempty"`
    Message string `json:"message"`
}

// 세션
type Session struct {
    Token string `json:"token"`
    PlaceID string `json:"place_id"`
    UserID string `json:"user_id"`
    ExpiredAt time.Time `json:"expired_at"`
}
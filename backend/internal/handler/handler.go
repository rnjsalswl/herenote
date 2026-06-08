package handler

import (
    "net/http"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/herenote/backend/internal/model"
    "github.com/herenote/backend/internal/service"
)

// 위치 인증
type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler(s *service.AuthService) *AuthHandler {
    return &AuthHandler{svc: s}
}

// POST /api/places/:id/verify
func (h *AuthHandler) Verify(c *gin.Context) {
    var req model.VerifyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.svc.Verify(c.Request.Context(), c.Param("id"), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    status := http.StatusOK
    if !resp.Verified {
        status = http.StatusForbidden
    }
    c.JSON(status, resp)
}

// 장소
type PlaceHandler struct{ svc *service.PlaceService }

func NewPlaceHandler(s *service.PlaceService) *PlaceHandler {
 return &PlaceHandler{svc: s}
}

// GET /api/places
func (h *PlaceHandler) List(c *gin.Context) {
    places, err := h.svc.List(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"places": places})
}

// GET /api/places/:id
func (h *PlaceHandler) Get(c *gin.Context) {
    place, err := h.svc.Get(c.Request.Context(), c.Param("id"))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "장소를 찾을 수 없습니다"})
        return
    }
    c.JSON(http.StatusOK, place)
}

// POST  /api/admin/places
func (h *PlaceHandler) Create(c *gin.Context) {
    var req model.CreatePlaceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    place, err := h.svc.Create(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, place)
}

// GET /api/places/nearby?lat=&lng={
func (h *PlaceHandler) Nearby(c *gin.Context) {
    lat := c.Query("lat")
    lng := c.Query("lng")

    if lat == "" || lng == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "lat, lng 파라미터가 필요합니다"})
        return
    }

    var latF, lngF float64
    fmt.Sscanf(lat, "%f", &latF)
    fmt.Sscanf(lng, "%f", &lngF)

    places, err := h.svc.FindNearby(c.Request.Context(), latF, lngF)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"places": places})
}

// 방명록
type GuestbookHandler struct { svc *service.GuestbookService }

func NewGuestbookHandler(s *service.GuestbookService) *GuestbookHandler {
    return &GuestbookHandler{svc: s}
}

// GET /api/places/:id/guestbooks
func (h *GuestbookHandler) List(c *gin.Context) {
    list, err := h.svc.List(c.Request.Context(), c.Param("id"))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"guestbooks": list})
}

// POST /api/places/:id/guestbooks
func (h *GuestbookHandler) Create(c *gin.Context) {
    var req model.CreateGuestbookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 미들웨어에서 주입된 값
    userID := c.GetString("user_id")
    badgeType := model.BadgeType(c.GetString("badgeType"))

    gb, err := h.svc.Create(c.Request.Context(), c.Param("id"), userID, &req, badgeType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gb)
}

// GET /api/users/:id/guestbooks/places
func (h *GuestbookHandler) MyPlaces(c *gin.Context) {
    places, err := h.svc.MyPlaces(c.Request.Context(), c.Param("id"))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"places": places})
}

// 유저
type UserHandler struct { svc *service.UserService }

func NewUserHandler(s *service.UserService) *UserHandler {
    return &UserHandler{svc: s}
}

// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
    var req model.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, err := h.svc.Create(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}

// GET /api/users/:id
func (h *UserHandler) Get(c *gin.Context) {
    user, err := h.svc.Get(c.Request.Context(), c.Param("id"))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "유저를 찾을 수 없습니다."})
        return
    }
    c.JSON(http.StatusOK, user)
}

// POST /api/admin/badges/:userId
func (h *UserHandler) GrantBadge(c *gin.Context) {
    var body struct {
        BadgeType string `json:"badge_type" binding:"required"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.svc.GrantBadge(c.Request.Context(), c.Param("userId"), model.BadgeType(body.BadgeType)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "뱃지가 부여되었습니다."})
}

// POST /api/auth/login
func (h *UserHandler) Login(c *gin.Context) {
    var req model.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, err := h.svc.Login(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, model.LoginResponse{
        UserID:    user.ID,
        Nickname:  user.Nickname,
        BadgeType: user.BadgeType,
    })
}
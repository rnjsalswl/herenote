package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/herenote/backend/internal/repository"
    "github.com/herenote/backend/internal/model"
)

// CORS - Vue PWA에서 API 호출 허용
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Location-token, X-Admin-Secret")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}

// LocationAuth - 방명록 열람, 작성 전 위치 인증 토큰 검증
func LocationAuth(sessionRepo *repository.SessionRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("X-Location-Token")
        placeID := c.Param("id")

        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "위치 인증이 필요합니다. 먼저 /verify를 호출해주세요.",
            })
            return
        }

        session, err := sessionRepo.FindValid(c.Request.Context(), token, placeID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "유효하지 않거나 만료된 인증입니다. 현장에서 다시 인증해주세요.",
            })
            return
        }

        // 핸들러에서 쓸 수 있게 컨텍스트에 주입
        c.Set("user_id", session.UserID)
        c.Set("place_id", session.PlaceID)
        c.Set("badge_type", string(model.BadgeNone))
        c.Next()
    }
}

// AdminAuth - 어드민 API 보호
func AdminAuth(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetHeader("X-Admin-Secret") != secret {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "관리자 권한이 없습니다.",
            })
            return
        }
        c.Next()
    }
}
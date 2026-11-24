package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    jwtauth "github.com/herenote/backend/internal/auth"
    "github.com/herenote/backend/internal/repository"
)

// CORS - Vue PWA에서 API 호출 허용
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Location-Token")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}

// UserAuth - JWT Bearer 토큰 검증 → user_id / nickname 컨텍스트 주입
func UserAuth(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if !strings.HasPrefix(header, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "인증이 필요합니다. Authorization: Bearer <token>",
            })
            return
        }
        tokenStr := strings.TrimPrefix(header, "Bearer ")
        claims, err := jwtauth.Verify(tokenStr, jwtSecret)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "유효하지 않거나 만료된 토큰입니다.",
            })
            return
        }
        c.Set("user_id", claims.UserID)
        c.Set("nickname", claims.Nickname)
        c.Next()
    }
}

// LocationAuth - 방명록 열람·작성 전 위치 인증 토큰 검증
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

        c.Set("user_id", session.UserID)
        c.Set("place_id", session.PlaceID)
        c.Next()
    }
}

// SelfOnly - 경로의 :param 이 로그인한 본인 ID일 때만 통과
// UserAuth 뒤에 붙여서 사용. 로그인만 확인하면 남의 id를 넣어서
// 다녀간 장소·방명록을 볼 수 있었음(IDOR). 방명록은 "현장에 가야만 읽기"가 규칙이라
// 이 경로로 위치 인증을 우회할 수 있었다.
func SelfOnly(param string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Param(param) == "" || c.Param(param) != c.GetString("user_id") {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "본인 정보만 조회할 수 있습니다."})
            return
        }
        c.Next()
    }
}

// AdminOnly - 어드민 API 보호
// UserAuth 뒤에 붙여서 사용. JWT로 확인된 user_id가 관리자 목록에 있어야 통과.
// (공유 시크릿 헤더 방식은 프론트 번들에 값이 노출되므로 사용하지 않음)
func AdminOnly(adminIDs map[string]bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !adminIDs[c.GetString("user_id")] {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "관리자 권한이 없습니다.",
            })
            return
        }
        c.Next()
    }
}

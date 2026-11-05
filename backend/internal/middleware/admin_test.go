package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwtauth "github.com/herenote/backend/internal/auth"
)

func TestAdminOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "0123456789abcdef0123456789abcdef"

	r := gin.New()
	r.POST("/admin", UserAuth(secret), AdminOnly(map[string]bool{"admin-1": true}), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	call := func(token string) int {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/admin", nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		// 예전 공유 시크릿 헤더는 더 이상 효과가 없어야 함
		req.Header.Set("X-Admin-Secret", "herenote-admin")
		r.ServeHTTP(w, req)
		return w.Code
	}

	adminToken, _ := jwtauth.Sign("admin-1", "admin", secret)
	userToken, _ := jwtauth.Sign("user-2", "user", secret)

	if got := call(""); got != http.StatusUnauthorized {
		t.Errorf("토큰 없음: want 401, got %d", got)
	}
	if got := call(userToken); got != http.StatusForbidden {
		t.Errorf("일반 유저: want 403, got %d", got)
	}
	if got := call(adminToken); got != http.StatusOK {
		t.Errorf("관리자: want 200, got %d", got)
	}
}

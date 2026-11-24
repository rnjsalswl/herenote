package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwtauth "github.com/herenote/backend/internal/auth"
)

func TestSelfOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "0123456789abcdef0123456789abcdef"

	r := gin.New()
	r.GET("/users/:id/places/:placeId/guestbooks", UserAuth(secret), SelfOnly("id"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	call := func(token, userID string) int {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/"+userID+"/places/p-1/guestbooks", nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		r.ServeHTTP(w, req)
		return w.Code
	}

	aliceToken, _ := jwtauth.Sign("alice", "alice", secret)

	if got := call("", "alice"); got != http.StatusUnauthorized {
		t.Errorf("토큰 없음: want 401, got %d", got)
	}
	if got := call(aliceToken, "bob"); got != http.StatusForbidden {
		t.Errorf("남의 id로 조회: want 403, got %d", got)
	}
	if got := call(aliceToken, "alice"); got != http.StatusOK {
		t.Errorf("본인 조회: want 200, got %d", got)
	}
}

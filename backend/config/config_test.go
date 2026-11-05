package config

import "testing"

func TestLoadRequiresJWTSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	if _, err := Load(); err == nil {
		t.Fatal("JWT_SECRET이 비어 있으면 에러가 나야 함")
	}

	t.Setenv("JWT_SECRET", "too-short")
	if _, err := Load(); err == nil {
		t.Fatal("JWT_SECRET이 32자 미만이면 에러가 나야 함")
	}
}

func TestLoadParsesAdminIDs(t *testing.T) {
	t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("ADMIN_USER_IDS", " a, ,b")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.AdminUserIDs) != 2 || !cfg.AdminUserIDs["a"] || !cfg.AdminUserIDs["b"] {
		t.Fatalf("unexpected admin ids: %v", cfg.AdminUserIDs)
	}
}

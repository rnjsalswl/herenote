package config

import (
    "fmt"
    "os"
    "strings"
)

type Config struct {
    Port         string
    DatabaseURL  string
    JWTSecret    string
    AdminUserIDs map[string]bool // 관리자 권한을 가진 유저 ID 목록
}

// Load는 환경변수에서 설정을 읽습니다.
// 보안에 필요한 값이 비어 있으면 기본값으로 뜨지 않고 에러를 반환합니다.
func Load() (*Config, error) {
    cfg := &Config{
        Port:         getEnv("PORT", "8088"),
        DatabaseURL:  os.Getenv("DATABASE_URL"),
        JWTSecret:    os.Getenv("JWT_SECRET"),
        AdminUserIDs: parseList(os.Getenv("ADMIN_USER_IDS")),
    }

    if len(cfg.JWTSecret) < 32 {
        return nil, fmt.Errorf("JWT_SECRET은 32자 이상으로 설정해야 합니다")
    }
    return cfg, nil
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}

// parseList는 "a,b,c" 형태의 값을 set으로 변환합니다.
func parseList(raw string) map[string]bool {
    set := make(map[string]bool)
    for _, v := range strings.Split(raw, ",") {
        if v = strings.TrimSpace(v); v != "" {
            set[v] = true
        }
    }
    return set
}

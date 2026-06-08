package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/herenote/backend/config"
    "github.com/herenote/backend/internal/handler"
    "github.com/herenote/backend/internal/middleware"
    "github.com/herenote/backend/internal/repository"
    "github.com/herenote/backend/internal/service"
    "github.com/joho/godotenv"
)

func main() {
    // .env 로드
    if err := godotenv.Load(); err != nil {
        log.Println(".env 파일 없음 - 환경변수 직접 사용")
    }

    // 설정 로드
    cfg := config.Load()

    // DB 연결
    db, err := repository.NewDB(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("DB 연결 실패: %v", err)
    }
    defer db.Close()

    // repository
    placeRepo := repository.NewPlaceRepository(db)
    guestbookRepo := repository.NewGuestbookRepository(db)
    sessionRepo := repository.NewSessionRepository(db)
    userRepo := repository.NewUserRepository(db)

    // service
    placeService := service.NewPlaceService(placeRepo)
    guestbookService := service.NewGuestbookService(guestbookRepo, sessionRepo)
    authService := service.NewAuthService(sessionRepo, placeRepo)
    userService := service.NewUserService(userRepo)

    // Handler
    placeHandler := handler.NewPlaceHandler(placeService)
    guestbookHandler := handler.NewGuestbookHandler(guestbookService)
    authHandler := handler.NewAuthHandler(authService)
    userHandler := handler.NewUserHandler(userService)

    // router
    r := gin.Default()
    r.Use(middleware.CORS())

    api := r.Group("/api")
    {
        // 위치 인증
        api.POST("/places/:id/verify", authHandler.Verify)

        // 장소
        api.GET("/places", placeHandler.List)
        api.GET("/places/:id", placeHandler.Get)
        api.GET("/places/nearby", placeHandler.Nearby)

        // 방명록 - 위치 인증 필요
        auth := api.Group("/places/:id/guestbooks")
        auth.Use(middleware.LocationAuth(sessionRepo))
        {
            auth.GET("", guestbookHandler.List)
            auth.POST("", guestbookHandler.Create)
        }
        api.GET("/users/:id/guestbooks/places", guestbookHandler.MyPlaces)

        // 유저
        api.POST("/users", userHandler.Create)
        api.GET("/users/:id", userHandler.Get)

        // 어드민
        admin := api.Group("/admin")
        admin.Use(middleware.AdminAuth(cfg.AdminSecret))
        {
            admin.POST("/places", placeHandler.Create)
            admin.POST("/badges/:userId", userHandler.GrantBadge)
        }
    }

    log.Printf("서버 시작: %s", cfg.Port)
    r.Run(":" + cfg.Port)
}

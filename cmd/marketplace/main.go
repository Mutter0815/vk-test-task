package main

import (
	"log"

	"github.com/Mutter0815/marketplace/internal/config"
	"github.com/Mutter0815/marketplace/internal/controllers"
	"github.com/Mutter0815/marketplace/internal/db"
	"github.com/Mutter0815/marketplace/internal/logger"
	"github.com/Mutter0815/marketplace/internal/middleware"
	"github.com/Mutter0815/marketplace/internal/repository"
	"github.com/Mutter0815/marketplace/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("не удалось загрузить конфиг: ", err)
	}
	lg, err := logger.New()
	if err != nil {
		log.Fatal("logger init: ", err)
	}
	logg := lg.Sugar()
	defer logg.Sync()
	logg.Info("logger initialized")

	db, err := db.Connect(cfg, logg)
	if err != nil {
		logg.Fatalw("Не удалось подключиться к базе данных", "err", err)
	}
	defer db.Close()

	repo := repository.NewSQLRepository(db, logg)
	authService := service.NewAuthService(repo, logg)
	adService := service.NewAdService(repo, logg)
	authCtrl := controllers.NewAuthController(authService, []byte(cfg.JWTSecret), logg)
	adsCtrl := controllers.NewADSController(adService, logg)

	router := gin.Default()

	router.POST("/auth/register", authCtrl.Register)
	router.POST("/auth/login", authCtrl.Login)
	router.GET("/ads", adsCtrl.ListAds)

	authGroup := router.Group("/")
	authGroup.Use(middleware.JWTAuthMiddleware([]byte(cfg.JWTSecret), logg))
	authGroup.POST("/ads", adsCtrl.CreateAd)
	logg.Info("listening on :8080")
	router.Run(":8080")

}

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yayodelmal/go-sensor-data-collector/internal/config"
	"github.com/yayodelmal/go-sensor-data-collector/internal/db"
	"github.com/yayodelmal/go-sensor-data-collector/internal/handlers"
	"github.com/yayodelmal/go-sensor-data-collector/internal/middleware"
	"github.com/yayodelmal/go-sensor-data-collector/internal/repository"
)

func main() {
	// 0) Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[main] config error: %v", err)
	}

	// 1) Conectar a TimescaleDB
	db.Init(cfg.DatabaseURL)

	// 2) Crear repositorios con inyección de dependencias
	sht31Repo := repository.NewSHT31Repository(db.DB)
	ds18b20Repo := repository.NewDS18B20Repository(db.DB)

	// 3) Crear handlers con repositorios inyectados
	sht31Handler := handlers.NewSHT31Handler(sht31Repo)
	ds18b20Handler := handlers.NewDS18B20Handler(ds18b20Repo)

	// 4) Crear router Gin
	router := gin.Default()

	// 5) Health check endpoint (sin autenticación)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// 6) Grupo de rutas bajo /sensor con autenticación opcional
	apiGroup := router.Group("/sensor")

	// Aplicar middleware solo al grupo /sensor si hay API_TOKEN
	if cfg.APIToken != "" {
		log.Println("[main] API_TOKEN detected → enabling auth middleware for /sensor routes")
		apiGroup.Use(middleware.AuthMiddleware())
	} else {
		log.Println("[main] No API_TOKEN found → /sensor endpoints are unprotected")
	}

	sht31Handler.RegisterRoutes(apiGroup)
	ds18b20Handler.RegisterRoutes(apiGroup)

	// 8) Arrancar servidor
	log.Printf("[main] starting server on :%s …", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("[main] failed to run server: %v", err)
	}
}

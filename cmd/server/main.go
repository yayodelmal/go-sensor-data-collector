package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yayodelmal/go-sensor-data-collector/internal/config"
	"github.com/yayodelmal/go-sensor-data-collector/internal/db"
	"github.com/yayodelmal/go-sensor-data-collector/internal/handlers"
	"github.com/yayodelmal/go-sensor-data-collector/internal/middleware"
)

func main() {
	// 0) Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[main] config error: %v", err)
	}

	// 1) Conectar a TimescaleDB
	db.Init(cfg.DatabaseURL)

	// 2) Crear router Gin
	router := gin.Default()

	// 3) Autenticación opcional
	if cfg.APIToken != "" {
		log.Println("[main] API_TOKEN detected → enabling auth middleware")
		router.Use(middleware.AuthMiddleware())
	} else {
		log.Println("[main] No API_TOKEN found → endpoints are unprotected")
	}

	// 4) Grupo de rutas bajo /sensor
	apiGroup := router.Group("/sensor")
	apiGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "go-sensor-data-collector is up"})
	})
	handlers.RegisterSHT31Routes(apiGroup)
	handlers.RegisterDS18B20Routes(apiGroup)

	// 5) Arrancar servidor
	log.Printf("[main] starting server on :%s …", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("[main] failed to run server: %v", err)
	}
}

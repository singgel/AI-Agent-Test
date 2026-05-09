package main

import (
	"fmt"
	"go-user-service/config"
	"go-user-service/models"
	"go-user-service/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	// 检查是否使用 Mock 模式
	useMock := os.Getenv("USE_MOCK") == "true"

	if useMock {
		log.Println("========================================")
		log.Println("Running in MOCK MODE")
		log.Println("Using in-memory data storage (no database required)")
		log.Println("========================================")
		models.InitMockDB()
		routes.SetupMockRoutes(r, cfg)
	} else {
		if err := models.InitDB(&cfg.Database); err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		routes.SetupRoutes(r, cfg)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

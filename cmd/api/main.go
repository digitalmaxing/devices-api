package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/digitalmaxing/devices-api/internal/config"
	"github.com/digitalmaxing/devices-api/internal/handlers"
	"github.com/digitalmaxing/devices-api/internal/models"
	"github.com/digitalmaxing/devices-api/internal/repository"
	"github.com/digitalmaxing/devices-api/internal/service"
)

func main() {
	cfg := config.Load()

	// Establish database connection
	db, err := gorm.Open(postgres.Open(cfg.GetDBDSN()), &gorm.Config{
		// Add any GORM config like logger here if needed
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate schema (for development; use proper migrations in production)
	if err := db.AutoMigrate(&models.Device{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Dependency injection: repo -> service -> handler
	repo := repository.NewPostgresDeviceRepository(db)
	svc := service.NewDeviceService(repo)
	deviceHandler := handlers.NewDeviceHandler(svc)

	// Initialize Gin router (release mode in prod via env)
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Health check for container orchestration / k8s
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "devices-api"})
	})

	// API routes group
	api := r.Group("/devices")
	{
		api.POST("", deviceHandler.CreateDevice)
		api.GET("", deviceHandler.ListDevices)
		api.GET("/:id", deviceHandler.GetDevice)
		api.PATCH("/:id", deviceHandler.UpdateDevice)
		api.DELETE("/:id", deviceHandler.DeleteDevice)
	}

	log.Printf("🚀 Devices API starting on %s | DB: %s:%s",
		cfg.GetServerAddr(), cfg.DBHost, cfg.DBPort)

	if err := r.Run(cfg.GetServerAddr()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
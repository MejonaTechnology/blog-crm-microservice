package main

import (
	"blog-service/internal/handlers"
	"blog-service/internal/middleware"
	"blog-service/pkg/database"
	"blog-service/pkg/logger"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var startTime = time.Now()

// @title Blog CRM Management Microservice API
// @version 1.0
// @description Professional blog and content management microservice for Mejona Technology CRM
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.mejona.com/support
// @contact.email support@mejona.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8082
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger
	logger.InitLogger()

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	
	// Add essential middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()

	// ===== HEALTH CHECK ENDPOINTS =====
	router.GET("/health", healthHandler.SimpleHealthCheck)
	router.GET("/health/deep", healthHandler.DeepHealthCheck)
	router.GET("/status", healthHandler.StatusCheck)
	router.GET("/ready", healthHandler.ReadinessCheck)
	router.GET("/alive", healthHandler.LivenessCheck)
	router.GET("/metrics", healthHandler.MetricsCheck)

	// ===== API ROUTES =====
	api := router.Group("/api/v1")
	{
		// Test endpoint
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Blog service test endpoint working",
				"data": map[string]interface{}{
					"service":   "Blog CRM Management Microservice",
					"version":   "1.0.0",
					"status":    "operational",
					"port":      "8082",
					"timestamp": time.Now(),
				},
			})
		})
	}

	// Swagger documentation
	if os.Getenv("ENABLE_SWAGGER") == "true" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Blog CRM Service starting on port %s", port)
	log.Printf("All endpoints available:")
	log.Printf("  HEALTH ENDPOINTS:")
	log.Printf("    GET  /health - Basic health check")
	log.Printf("    GET  /health/deep - Deep health check")
	log.Printf("    GET  /status - Quick status")
	log.Printf("    GET  /ready - Readiness check")
	log.Printf("    GET  /alive - Liveness check")
	log.Printf("    GET  /metrics - System metrics")
	log.Printf("  API ENDPOINTS:")
	log.Printf("    GET  /api/v1/test - Test endpoint")
	log.Printf("  DOCUMENTATION:")
	log.Printf("    GET  /swagger/index.html - API Documentation (if enabled)")
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
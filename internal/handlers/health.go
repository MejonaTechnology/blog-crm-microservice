package handlers

import (
	"blog-service/pkg/database"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// SimpleHealthCheck performs a basic health check
func (h *HealthHandler) SimpleHealthCheck(c *gin.Context) {
	uptime := time.Since(startTime)
	
	// Check database
	dbHealth := "healthy"
	if db := database.GetDB(); db != nil {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Ping(); err != nil {
				dbHealth = "unhealthy"
			}
		} else {
			dbHealth = "unhealthy"
		}
	} else {
		dbHealth = "unhealthy"
	}

	status := "healthy"
	if dbHealth != "healthy" {
		status = "unhealthy"
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Health check completed",
		"data": map[string]interface{}{
			"status":    status,
			"service":   "Blog CRM Management Microservice",
			"version":   getEnv("APP_VERSION", "1.0.0"),
			"uptime":    uptime.String(),
			"database":  dbHealth,
			"port":      getEnv("PORT", "8082"),
			"timestamp": time.Now(),
		},
	}

	statusCode := http.StatusOK
	if status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// DeepHealthCheck performs comprehensive health checks
func (h *HealthHandler) DeepHealthCheck(c *gin.Context) {
	startCheck := time.Now()
	checks := make(map[string]interface{})
	overallStatus := "healthy"

	// Database check
	dbStart := time.Now()
	dbHealth := "healthy"
	var dbError error
	if db := database.GetDB(); db != nil {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Ping(); err != nil {
				dbHealth = "unhealthy"
				dbError = err
				overallStatus = "unhealthy"
			} else {
				// Test a simple query
				if err := database.TestQuery(); err != nil {
					dbHealth = "degraded"
					dbError = err
					if overallStatus == "healthy" {
						overallStatus = "degraded"
					}
				}
			}
		} else {
			dbHealth = "unhealthy"
			dbError = err
			overallStatus = "unhealthy"
		}
	} else {
		dbHealth = "unhealthy"
		overallStatus = "unhealthy"
	}
	dbDuration := time.Since(dbStart)

	checks["database"] = map[string]interface{}{
		"status":      dbHealth,
		"duration_ms": dbDuration.Milliseconds(),
		"error":       dbError,
		"stats":       database.GetConnectionStats(),
	}

	// Memory check
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryStatus := "healthy"
	allocatedMB := memStats.Alloc / 1024 / 1024

	if allocatedMB > 1024 { // 1GB
		memoryStatus = "critical"
		overallStatus = "critical"
	} else if allocatedMB > 512 { // 512MB
		memoryStatus = "warning"
		if overallStatus == "healthy" {
			overallStatus = "warning"
		}
	}

	checks["memory"] = map[string]interface{}{
		"status":         memoryStatus,
		"allocated_mb":   allocatedMB,
		"heap_in_use_mb": memStats.HeapInuse / 1024 / 1024,
		"gc_count":       memStats.NumGC,
		"gc_pause_ns":    memStats.PauseTotalNs,
	}

	// Disk space check (optional)
	checks["disk"] = map[string]interface{}{
		"status": "healthy",
		"note":   "Disk space monitoring not implemented",
	}

	// Environment check
	envStatus := "healthy"
	missingVars := []string{}
	
	requiredVars := []string{"DB_HOST", "DB_NAME", "JWT_SECRET"}
	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
			envStatus = "degraded"
			if overallStatus == "healthy" {
				overallStatus = "degraded"
			}
		}
	}

	checks["environment"] = map[string]interface{}{
		"status":       envStatus,
		"missing_vars": missingVars,
	}

	totalDuration := time.Since(startCheck)
	statusCode := http.StatusOK
	if overallStatus == "critical" || overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	} else if overallStatus == "degraded" || overallStatus == "warning" {
		statusCode = http.StatusOK // Still operational but with issues
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Deep health check completed",
		"data": map[string]interface{}{
			"status":             overallStatus,
			"check_duration_ms":  totalDuration.Milliseconds(),
			"checks":             checks,
			"service":            "Blog CRM Management Microservice",
			"version":            getEnv("APP_VERSION", "1.0.0"),
			"uptime":             time.Since(startTime).String(),
			"timestamp":          time.Now(),
		},
	}

	c.JSON(statusCode, response)
}

// StatusCheck returns quick status information
func (h *HealthHandler) StatusCheck(c *gin.Context) {
	response := map[string]interface{}{
		"success": true,
		"message": "Status retrieved",
		"data": map[string]interface{}{
			"status":         "healthy",
			"uptime_seconds": time.Since(startTime).Seconds(),
			"version":        getEnv("APP_VERSION", "1.0.0"),
			"environment":    getEnv("APP_ENV", "production"),
			"service":        "Blog CRM Management Microservice",
			"port":           getEnv("PORT", "8082"),
			"timestamp":      time.Now(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// ReadinessCheck checks if service is ready to serve requests
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Check if database is connected and responsive
	ready := true
	reasons := []string{}

	if db := database.GetDB(); db != nil {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Ping(); err != nil {
				ready = false
				reasons = append(reasons, "database connection failed")
			}
		} else {
			ready = false
			reasons = append(reasons, "unable to get database instance")
		}
	} else {
		ready = false
		reasons = append(reasons, "database not initialized")
	}

	// Check essential environment variables
	requiredEnvVars := []string{"DB_HOST", "DB_NAME"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			ready = false
			reasons = append(reasons, "missing required environment variable: "+envVar)
		}
	}

	statusCode := http.StatusOK
	message := "Service is ready"
	if !ready {
		statusCode = http.StatusServiceUnavailable
		message = "Service not ready"
	}

	response := map[string]interface{}{
		"success": ready,
		"message": message,
		"data": map[string]interface{}{
			"ready":     ready,
			"reasons":   reasons,
			"timestamp": time.Now(),
		},
	}

	c.JSON(statusCode, response)
}

// LivenessCheck checks if service is alive and responsive
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	response := map[string]interface{}{
		"success": true,
		"message": "Service is alive",
		"data": map[string]interface{}{
			"alive":     true,
			"service":   "Blog CRM Management Microservice",
			"uptime":    time.Since(startTime).String(),
			"timestamp": time.Now(),
		},
	}

	c.JSON(http.StatusOK, response)
}

// MetricsCheck returns comprehensive system metrics
func (h *HealthHandler) MetricsCheck(c *gin.Context) {
	uptime := time.Since(startTime)
	
	// Get memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := map[string]interface{}{
		"service": map[string]interface{}{
			"name":           "Blog CRM Management Microservice",
			"uptime_seconds": uptime.Seconds(),
			"version":        getEnv("APP_VERSION", "1.0.0"),
			"environment":    getEnv("APP_ENV", "production"),
			"port":           getEnv("PORT", "8082"),
			"start_time":     startTime,
		},
		"runtime": map[string]interface{}{
			"go_version":   runtime.Version(),
			"go_routines":  runtime.NumGoroutine(),
			"go_max_procs": runtime.GOMAXPROCS(0),
			"memory": map[string]interface{}{
				"allocated_mb":       memStats.Alloc / 1024 / 1024,
				"total_allocated_mb": memStats.TotalAlloc / 1024 / 1024,
				"system_mb":          memStats.Sys / 1024 / 1024,
				"heap_allocated_mb":  memStats.HeapAlloc / 1024 / 1024,
				"heap_in_use_mb":     memStats.HeapInuse / 1024 / 1024,
				"gc_count":           memStats.NumGC,
				"gc_pause_total_ns":  memStats.PauseTotalNs,
			},
		},
		"database": database.GetConnectionStats(),
		"timestamp": time.Now(),
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Comprehensive metrics retrieved",
		"data":    metrics,
	}

	c.JSON(http.StatusOK, response)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
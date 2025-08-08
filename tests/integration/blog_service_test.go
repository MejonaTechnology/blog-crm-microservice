// +build integration

package integration

import (
	"blog-service/internal/handlers"
	"blog-service/internal/middleware"
	"blog-service/pkg/database"
	"blog-service/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// BlogServiceIntegrationSuite contains integration tests for the blog service
type BlogServiceIntegrationSuite struct {
	suite.Suite
	router *gin.Engine
	server *httptest.Server
}

// SetupSuite initializes the test suite
func (suite *BlogServiceIntegrationSuite) SetupSuite() {
	// Load test environment
	if err := godotenv.Load("../../.env.test"); err != nil {
		suite.T().Logf("No .env.test file found, using system environment")
	}

	// Set test mode
	gin.SetMode(gin.TestMode)
	os.Setenv("GIN_MODE", "test")

	// Initialize logger for testing
	logger.InitLogger()

	// Initialize test database
	err := database.InitDB()
	suite.Require().NoError(err, "Failed to initialize test database")

	// Setup router
	suite.setupRouter()

	// Create test server
	suite.server = httptest.NewServer(suite.router)
}

// TearDownSuite cleans up after all tests
func (suite *BlogServiceIntegrationSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Close()
	}
	
	// Cleanup test database if needed
	// Note: Add your database cleanup logic here
}

// setupRouter configures the Gin router for testing
func (suite *BlogServiceIntegrationSuite) setupRouter() {
	suite.router = gin.New()
	
	// Add middleware
	suite.router.Use(gin.Logger())
	suite.router.Use(gin.Recovery())
	suite.router.Use(middleware.CORS())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()

	// Health check endpoints
	suite.router.GET("/health", healthHandler.SimpleHealthCheck)
	suite.router.GET("/health/deep", healthHandler.DeepHealthCheck)
	suite.router.GET("/status", healthHandler.StatusCheck)
	suite.router.GET("/ready", healthHandler.ReadinessCheck)
	suite.router.GET("/alive", healthHandler.LivenessCheck)
	suite.router.GET("/metrics", healthHandler.MetricsCheck)

	// API routes
	api := suite.router.Group("/api/v1")
	{
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
}

// TestHealthEndpoints tests all health check endpoints
func (suite *BlogServiceIntegrationSuite) TestHealthEndpoints() {
	testCases := []struct {
		name           string
		endpoint       string
		expectedStatus int
		checkResponse  bool
	}{
		{"Simple Health Check", "/health", http.StatusOK, true},
		{"Deep Health Check", "/health/deep", http.StatusOK, true},
		{"Status Check", "/status", http.StatusOK, true},
		{"Readiness Check", "/ready", http.StatusOK, true},
		{"Liveness Check", "/alive", http.StatusOK, true},
		{"Metrics Check", "/metrics", http.StatusOK, false},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := http.Get(suite.server.URL + tc.endpoint)
			suite.Require().NoError(err)
			defer resp.Body.Close()

			suite.Equal(tc.expectedStatus, resp.StatusCode, 
				"Expected status %d for %s, got %d", tc.expectedStatus, tc.endpoint, resp.StatusCode)

			if tc.checkResponse {
				var result map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&result)
				suite.Require().NoError(err)
				suite.Contains(result, "status", "Response should contain status field")
			}
		})
	}
}

// TestAPIEndpoints tests the main API endpoints
func (suite *BlogServiceIntegrationSuite) TestAPIEndpoints() {
	suite.Run("Test API Endpoint", func() {
		resp, err := http.Get(suite.server.URL + "/api/v1/test")
		suite.Require().NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		suite.Require().NoError(err)

		suite.Equal(true, result["success"])
		suite.Contains(result, "data")
		
		data, ok := result["data"].(map[string]interface{})
		suite.True(ok, "Data should be a map")
		suite.Equal("Blog CRM Management Microservice", data["service"])
		suite.Equal("1.0.0", data["version"])
		suite.Equal("operational", data["status"])
		suite.Equal("8082", data["port"])
	})
}

// TestCORSHeaders tests CORS middleware functionality
func (suite *BlogServiceIntegrationSuite) TestCORSHeaders() {
	suite.Run("CORS Headers Present", func() {
		req, err := http.NewRequest("OPTIONS", suite.server.URL+"/api/v1/test", nil)
		suite.Require().NoError(err)
		
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")

		client := &http.Client{}
		resp, err := client.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusNoContent, resp.StatusCode)
		suite.Contains(resp.Header.Get("Access-Control-Allow-Origin"), "*")
		suite.Contains(resp.Header.Get("Access-Control-Allow-Methods"), "GET")
	})
}

// TestResponseTimes tests API response time performance
func (suite *BlogServiceIntegrationSuite) TestResponseTimes() {
	endpoints := []string{"/health", "/api/v1/test", "/status"}
	maxResponseTime := 1000 * time.Millisecond // 1 second

	for _, endpoint := range endpoints {
		suite.Run(fmt.Sprintf("Response Time %s", endpoint), func() {
			start := time.Now()
			
			resp, err := http.Get(suite.server.URL + endpoint)
			suite.Require().NoError(err)
			defer resp.Body.Close()

			elapsed := time.Since(start)
			suite.True(elapsed < maxResponseTime, 
				"Response time for %s took %v, expected less than %v", 
				endpoint, elapsed, maxResponseTime)
			
			suite.Equal(http.StatusOK, resp.StatusCode)
		})
	}
}

// TestConcurrentRequests tests service behavior under concurrent load
func (suite *BlogServiceIntegrationSuite) TestConcurrentRequests() {
	suite.Run("Concurrent Health Checks", func() {
		concurrency := 10
		requests := 50
		results := make(chan error, concurrency*requests)

		// Launch concurrent goroutines
		for i := 0; i < concurrency; i++ {
			go func() {
				for j := 0; j < requests; j++ {
					resp, err := http.Get(suite.server.URL + "/health")
					if err != nil {
						results <- err
						continue
					}
					
					if resp.StatusCode != http.StatusOK {
						results <- fmt.Errorf("expected status 200, got %d", resp.StatusCode)
						resp.Body.Close()
						continue
					}
					
					resp.Body.Close()
					results <- nil
				}
			}()
		}

		// Collect results
		successCount := 0
		errorCount := 0
		
		for i := 0; i < concurrency*requests; i++ {
			if err := <-results; err != nil {
				suite.T().Logf("Request failed: %v", err)
				errorCount++
			} else {
				successCount++
			}
		}

		suite.True(successCount > errorCount, 
			"Expected more successes than failures. Success: %d, Errors: %d", 
			successCount, errorCount)
		
		suite.True(float64(successCount)/float64(successCount+errorCount) > 0.95,
			"Success rate should be above 95%")
	})
}

// TestErrorHandling tests error scenarios
func (suite *BlogServiceIntegrationSuite) TestErrorHandling() {
	suite.Run("404 Not Found", func() {
		resp, err := http.Get(suite.server.URL + "/nonexistent")
		suite.Require().NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusNotFound, resp.StatusCode)
	})

	suite.Run("Invalid HTTP Method", func() {
		req, err := http.NewRequest("DELETE", suite.server.URL+"/health", nil)
		suite.Require().NoError(err)

		client := &http.Client{}
		resp, err := client.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusMethodNotAllowed, resp.StatusCode)
	})
}

// TestServiceRecovery tests service recovery after panic
func (suite *BlogServiceIntegrationSuite) TestServiceRecovery() {
	suite.Run("Service Recovery", func() {
		// Add a route that panics
		suite.router.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})

		resp, err := http.Get(suite.server.URL + "/panic")
		suite.Require().NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusInternalServerError, resp.StatusCode)

		// Verify service is still functional after panic
		resp, err = http.Get(suite.server.URL + "/health")
		suite.Require().NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusOK, resp.StatusCode)
	})
}

// TestJSONResponses tests JSON response format consistency
func (suite *BlogServiceIntegrationSuite) TestJSONResponses() {
	endpoints := []string{"/health", "/health/deep", "/api/v1/test"}

	for _, endpoint := range endpoints {
		suite.Run(fmt.Sprintf("JSON Response %s", endpoint), func() {
			resp, err := http.Get(suite.server.URL + endpoint)
			suite.Require().NoError(err)
			defer resp.Body.Close()

			suite.Equal("application/json; charset=utf-8", resp.Header.Get("Content-Type"))

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			suite.Require().NoError(err, "Response should be valid JSON")
		})
	}
}

// TestLargePayloads tests handling of large request payloads
func (suite *BlogServiceIntegrationSuite) TestLargePayloads() {
	suite.Run("Large POST Payload", func() {
		// Create a large payload (1MB)
		largeData := make(map[string]string)
		largeData["data"] = string(make([]byte, 1024*1024))

		payload, err := json.Marshal(largeData)
		suite.Require().NoError(err)

		// Add a test route that accepts POST
		suite.router.POST("/test/large", func(c *gin.Context) {
			var data map[string]string
			if err := c.ShouldBindJSON(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"received": len(data["data"])})
		})

		resp, err := http.Post(suite.server.URL+"/test/large", 
			"application/json", bytes.NewReader(payload))
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Should handle large payload gracefully
		suite.True(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusRequestEntityTooLarge)
	})
}

// TestDatabaseConnection tests database connectivity (if applicable)
func (suite *BlogServiceIntegrationSuite) TestDatabaseConnection() {
	suite.Run("Database Health Check", func() {
		resp, err := http.Get(suite.server.URL + "/health/deep")
		suite.Require().NoError(err)
		defer resp.Body.Close()

		suite.Equal(http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		suite.Require().NoError(err)

		// Check if database status is included in deep health check
		if database, exists := result["database"]; exists {
			dbStatus, ok := database.(map[string]interface{})
			suite.True(ok, "Database status should be a map")
			suite.Contains(dbStatus, "status", "Database status should include status field")
		}
	})
}

// TestSecurityHeaders tests security-related HTTP headers
func (suite *BlogServiceIntegrationSuite) TestSecurityHeaders() {
	suite.Run("Security Headers", func() {
		resp, err := http.Get(suite.server.URL + "/health")
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Check for security headers set by middleware
		headers := resp.Header

		// Check CORS headers
		suite.Contains(headers.Get("Access-Control-Allow-Origin"), "*")
		
		// Note: Additional security headers would be set by nginx in production
		// These tests verify application-level headers
	})
}

// Run the integration test suite
func TestBlogServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(BlogServiceIntegrationSuite))
}
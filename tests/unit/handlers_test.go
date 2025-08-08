package unit

import (
	"blog-service/internal/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HandlersTestSuite contains unit tests for handlers
type HandlersTestSuite struct {
	suite.Suite
	router *gin.Engine
}

// SetupSuite initializes the test suite
func (suite *HandlersTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	
	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()

	// Setup routes
	suite.router.GET("/health", healthHandler.SimpleHealthCheck)
	suite.router.GET("/health/deep", healthHandler.DeepHealthCheck)
	suite.router.GET("/status", healthHandler.StatusCheck)
	suite.router.GET("/ready", healthHandler.ReadinessCheck)
	suite.router.GET("/alive", healthHandler.LivenessCheck)
	suite.router.GET("/metrics", healthHandler.MetricsCheck)
}

// TestSimpleHealthCheck tests the basic health check endpoint
func (suite *HandlersTestSuite) TestSimpleHealthCheck() {
	suite.Run("Simple Health Check Success", func() {
		req, err := http.NewRequest("GET", "/health", nil)
		suite.Require().NoError(err)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		suite.Equal(http.StatusOK, w.Code)
		suite.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		suite.Require().NoError(err)

		// Check response structure
		suite.Contains(response, "status")
		suite.Contains(response, "timestamp")
		suite.Contains(response, "service")
		suite.Contains(response, "version")

		// Check values
		suite.Equal("healthy", response["status"])
		suite.Equal("blog-service", response["service"])
		suite.NotEmpty(response["timestamp"])
	})
}

// TestDeepHealthCheck tests the comprehensive health check endpoint
func (suite *HandlersTestSuite) TestDeepHealthCheck() {
	suite.Run("Deep Health Check Success", func() {
		req, err := http.NewRequest("GET", "/health/deep", nil)
		suite.Require().NoError(err)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		suite.Equal(http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		suite.Require().NoError(err)

		// Check comprehensive health response structure
		suite.Contains(response, "status")
		suite.Contains(response, "checks")
		suite.Contains(response, "timestamp")
		suite.Contains(response, "uptime")
		suite.Contains(response, "system")

		// Check that it includes system information
		system, ok := response["system"].(map[string]interface{})
		suite.True(ok, "System should be a map")
		suite.Contains(system, "memory")
		suite.Contains(system, "goroutines")
	})
}

// TestStatusCheck tests the status endpoint
func (suite *HandlersTestSuite) TestStatusCheck() {
	suite.Run("Status Check Success", func() {
		req, err := http.NewRequest("GET", "/status", nil)
		suite.Require().NoError(err)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		suite.Equal(http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		suite.Require().NoError(err)

		suite.Contains(response, "status")
		suite.Equal("OK", response["status"])
	})
}

// TestReadinessCheck tests the readiness probe endpoint
func (suite *HandlersTestSuite) TestReadinessCheck() {
	suite.Run("Readiness Check Success", func() {
		req, err := http.NewRequest("GET", "/ready", nil)
		suite.Require().NoError(err)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		suite.Equal(http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		suite.Require().NoError(err)

		suite.Contains(response, "ready")
		suite.Equal(true, response["ready"])
	})
}

// TestLivenessCheck tests the liveness probe endpoint
func (suite *HandlersTestSuite) TestLivenessCheck() {
	suite.Run("Liveness Check Success", func() {
		req, err := http.NewRequest("GET", "/alive", nil)
		suite.Require().NoError(err)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		suite.Equal(http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		suite.Require().NoError(err)

		suite.Contains(response, "alive")
		suite.Equal(true, response["alive"])
	})
}

// TestMetricsCheck tests the metrics endpoint
func (suite *HandlersTestSuite) TestMetricsCheck() {
	suite.Run("Metrics Check Success", func() {
		req, err := http.NewRequest("GET", "/metrics", nil)
		suite.Require().NoError(err)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		suite.Equal(http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		suite.Require().NoError(err)

		// Check metrics response structure
		suite.Contains(response, "metrics")
		suite.Contains(response, "timestamp")

		metrics, ok := response["metrics"].(map[string]interface{})
		suite.True(ok, "Metrics should be a map")
		suite.Contains(metrics, "requests_total")
		suite.Contains(metrics, "uptime_seconds")
	})
}

// TestHTTPMethods tests that endpoints only accept appropriate HTTP methods
func (suite *HandlersTestSuite) TestHTTPMethods() {
	endpoints := []string{"/health", "/health/deep", "/status", "/ready", "/alive", "/metrics"}
	
	for _, endpoint := range endpoints {
		suite.Run("POST Not Allowed on "+endpoint, func() {
			req, err := http.NewRequest("POST", endpoint, nil)
			suite.Require().NoError(err)

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			suite.Equal(http.StatusMethodNotAllowed, w.Code)
		})

		suite.Run("PUT Not Allowed on "+endpoint, func() {
			req, err := http.NewRequest("PUT", endpoint, nil)
			suite.Require().NoError(err)

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			suite.Equal(http.StatusMethodNotAllowed, w.Code)
		})

		suite.Run("DELETE Not Allowed on "+endpoint, func() {
			req, err := http.NewRequest("DELETE", endpoint, nil)
			suite.Require().NoError(err)

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			suite.Equal(http.StatusMethodNotAllowed, w.Code)
		})
	}
}

// TestResponseHeaders tests that appropriate headers are set
func (suite *HandlersTestSuite) TestResponseHeaders() {
	suite.Run("Content-Type Headers", func() {
		req, err := http.NewRequest("GET", "/health", nil)
		suite.Require().NoError(err)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		suite.Equal(http.StatusOK, w.Code)
		suite.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})
}

// TestResponseFormat tests JSON response format consistency
func (suite *HandlersTestSuite) TestResponseFormat() {
	testCases := []struct {
		endpoint     string
		requiredKeys []string
	}{
		{"/health", []string{"status", "timestamp", "service"}},
		{"/health/deep", []string{"status", "checks", "timestamp", "uptime"}},
		{"/status", []string{"status"}},
		{"/ready", []string{"ready"}},
		{"/alive", []string{"alive"}},
		{"/metrics", []string{"metrics", "timestamp"}},
	}

	for _, tc := range testCases {
		suite.Run("Response Format "+tc.endpoint, func() {
			req, err := http.NewRequest("GET", tc.endpoint, nil)
			suite.Require().NoError(err)

			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			suite.Equal(http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			suite.Require().NoError(err)

			for _, key := range tc.requiredKeys {
				suite.Contains(response, key, "Response should contain key: %s", key)
			}
		})
	}
}

// TestConcurrentHandlers tests handler behavior under concurrent requests
func (suite *HandlersTestSuite) TestConcurrentHandlers() {
	suite.Run("Concurrent Health Check Requests", func() {
		concurrency := 10
		results := make(chan int, concurrency)

		// Launch concurrent requests
		for i := 0; i < concurrency; i++ {
			go func() {
				req, err := http.NewRequest("GET", "/health", nil)
				if err != nil {
					results <- 500
					return
				}

				w := httptest.NewRecorder()
				suite.router.ServeHTTP(w, req)
				results <- w.Code
			}()
		}

		// Collect results
		successCount := 0
		for i := 0; i < concurrency; i++ {
			statusCode := <-results
			if statusCode == http.StatusOK {
				successCount++
			}
		}

		suite.Equal(concurrency, successCount, "All concurrent requests should succeed")
	})
}

// TestResponseTime tests that handlers respond quickly
func (suite *HandlersTestSuite) TestResponseTime() {
	endpoints := []string{"/health", "/status", "/ready", "/alive"}

	for _, endpoint := range endpoints {
		suite.Run("Response Time "+endpoint, func() {
			req, err := http.NewRequest("GET", endpoint, nil)
			suite.Require().NoError(err)

			w := httptest.NewRecorder()
			
			// Measure response time
			start := time.Now()
			suite.router.ServeHTTP(w, req)
			elapsed := time.Since(start)

			suite.Equal(http.StatusOK, w.Code)
			suite.True(elapsed < time.Millisecond*100, 
				"Handler %s should respond in less than 100ms, took %v", endpoint, elapsed)
		})
	}
}

// TestErrorConditions tests handler behavior in error scenarios
func (suite *HandlersTestSuite) TestErrorConditions() {
	suite.Run("Malformed Request", func() {
		req, err := http.NewRequest("GET", "/health", nil)
		suite.Require().NoError(err)
		
		// Add malformed headers
		req.Header.Set("Content-Type", "invalid/content/type")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		// Should still respond successfully as health check doesn't depend on headers
		suite.Equal(http.StatusOK, w.Code)
	})
}

// Run the unit test suite
func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))

	// Run additional standalone tests
	t.Run("Health Handler Creation", func(t *testing.T) {
		handler := handlers.NewHealthHandler()
		assert.NotNil(t, handler)
	})
}
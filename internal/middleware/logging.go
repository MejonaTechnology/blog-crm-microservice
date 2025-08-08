package middleware

import (
	"blog-service/pkg/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger middleware logs all HTTP requests with detailed information
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Extract user ID from context if available
		var userID *uint
		if param.Keys != nil {
			if uid, exists := param.Keys["user_id"]; exists {
				if id, ok := uid.(uint); ok {
					userID = &id
				}
			}
		}

		// Log the API request
		logger.LogAPIRequest(
			param.Method,
			param.Path,
			userID,
			param.Latency,
			param.StatusCode,
		)

		// Return empty string as we handle logging internally
		return ""
	})
}

// DetailedRequestLogger provides comprehensive request/response logging
func DetailedRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Start time
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Extract user ID from context
		var userID *uint
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uint); ok {
				userID = &id
			}
		}

		// Log detailed request information
		fields := map[string]interface{}{
			"request_id":     requestID,
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"query":          c.Request.URL.RawQuery,
			"status_code":    c.Writer.Status(),
			"duration_ms":    duration.Milliseconds(),
			"client_ip":      c.ClientIP(),
			"user_agent":     c.Request.UserAgent(),
			"referer":        c.Request.Referer(),
			"content_type":   c.Request.Header.Get("Content-Type"),
			"content_length": c.Request.ContentLength,
		}

		if userID != nil {
			fields["user_id"] = *userID
		}

		// Log based on status code
		if c.Writer.Status() >= 500 {
			logger.Error("HTTP request failed with server error", nil, fields)
		} else if c.Writer.Status() >= 400 {
			logger.Warn("HTTP request failed with client error", fields)
		} else if duration > 5*time.Second {
			logger.Warn("Slow HTTP request detected", fields)
		} else {
			logger.Info("HTTP request completed", fields)
		}

		// Log performance metrics
		logger.LogPerformanceMetric(
			"http_request_duration",
			float64(duration.Milliseconds()),
			"ms",
			map[string]string{
				"method":   c.Request.Method,
				"endpoint": c.FullPath(),
				"status":   strconv.Itoa(c.Writer.Status()),
			},
		)
	}
}

// ErrorLoggingMiddleware logs errors with comprehensive context
func ErrorLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			requestID, _ := c.Get("request_id")
			userID, _ := c.Get("user_id")

			for _, ginErr := range c.Errors {
				fields := map[string]interface{}{
					"request_id": requestID,
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"status":     c.Writer.Status(),
					"client_ip":  c.ClientIP(),
					"user_agent": c.Request.UserAgent(),
				}

				if userID != nil {
					fields["user_id"] = userID
				}

				logger.Error("Request processing error", ginErr.Err, fields)
			}
		}
	}
}

// PanicRecoveryMiddleware recovers from panics and logs them
func PanicRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID, _ := c.Get("request_id")
		userID, _ := c.Get("user_id")

		fields := map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"panic":      recovered,
		}

		if userID != nil {
			fields["user_id"] = userID
		}

		logger.Error("PANIC RECOVERED", nil, fields)

		// Log as security event as well (panics might indicate attacks)
		var uid *uint
		if userID != nil {
			if id, ok := userID.(uint); ok {
				uid = &id
			}
		}
		logger.LogSecurityEvent(
			"panic_recovered",
			uid,
			c.ClientIP(),
			map[string]interface{}{
				"panic_value": recovered,
				"endpoint":    c.Request.URL.Path,
				"method":      c.Request.Method,
			},
		)

		// Return error response
		c.JSON(500, gin.H{
			"success": false,
			"message": "Internal server error",
			"data":    nil,
		})
	})
}

// SecurityEventLogger logs security-related events
func SecurityEventLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log security events based on response
		userID, _ := c.Get("user_id")
		var uid *uint
		if userID != nil {
			if id, ok := userID.(uint); ok {
				uid = &id
			}
		}

		// Failed authentication attempts
		if c.Writer.Status() == 401 {
			logger.LogSecurityEvent(
				"authentication_failed",
				uid,
				c.ClientIP(),
				map[string]interface{}{
					"endpoint":   c.Request.URL.Path,
					"method":     c.Request.Method,
					"user_agent": c.Request.UserAgent(),
				},
			)
		}

		// Access denied events
		if c.Writer.Status() == 403 {
			logger.LogSecurityEvent(
				"access_denied",
				uid,
				c.ClientIP(),
				map[string]interface{}{
					"endpoint":   c.Request.URL.Path,
					"method":     c.Request.Method,
					"user_agent": c.Request.UserAgent(),
				},
			)
		}
	}
}

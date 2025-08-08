package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

// InitLogger initializes the logger with file rotation
func InitLogger() {
	Logger = logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// Set log format
	if getEnv("LOG_FORMAT", "json") == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// Set up file logging if enabled
	if getEnv("LOG_FILE_ENABLED", "true") == "true" {
		logPath := getEnv("LOG_FILE_PATH", "./logs/blog-service.log")

		// Ensure log directory exists
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			Logger.Warnf("Failed to create log directory: %v", err)
		}

		// Set up log rotation
		rotateLogger := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    100, // megabytes
			MaxBackups: 10,
			MaxAge:     30, // days
			Compress:   true,
		}

		// Write to both file and stdout
		Logger.SetOutput(io.MultiWriter(os.Stdout, rotateLogger))
	} else {
		Logger.SetOutput(os.Stdout)
	}

	Logger.Info("Logger initialized successfully")
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}

// Info logs an info message
func Info(message string, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.WithFields(fields).Info(message)
}

// Warn logs a warning message
func Warn(message string, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.WithFields(fields).Warn(message)
}

// Error logs an error message
func Error(message string, err error, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}
	
	if fields == nil {
		fields = make(map[string]interface{})
	}
	
	if err != nil {
		fields["error"] = err.Error()
	}
	
	Logger.WithFields(fields).Error(message)
}

// Debug logs a debug message
func Debug(message string, fields map[string]interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.WithFields(fields).Debug(message)
}

// LogAPIRequest logs API request details
func LogAPIRequest(method, path string, userID *uint, duration time.Duration, statusCode int) {
	fields := map[string]interface{}{
		"method":      method,
		"path":        path,
		"duration_ms": duration.Milliseconds(),
		"status_code": statusCode,
	}
	
	if userID != nil {
		fields["user_id"] = *userID
	}
	
	if statusCode >= 500 {
		Error("API request failed with server error", nil, fields)
	} else if statusCode >= 400 {
		Warn("API request failed with client error", fields)
	} else {
		Info("API request completed", fields)
	}
}

// LogPerformanceMetric logs performance metrics
func LogPerformanceMetric(metric string, value float64, unit string, tags map[string]string) {
	fields := map[string]interface{}{
		"metric": metric,
		"value":  value,
		"unit":   unit,
	}
	
	for k, v := range tags {
		fields["tag_"+k] = v
	}
	
	Debug("Performance metric recorded", fields)
}

// LogSecurityEvent logs security-related events
func LogSecurityEvent(event string, userID *uint, ipAddress string, details map[string]interface{}) {
	fields := map[string]interface{}{
		"security_event": event,
		"ip_address":     ipAddress,
	}
	
	if userID != nil {
		fields["user_id"] = *userID
	}
	
	for k, v := range details {
		fields[k] = v
	}
	
	Warn("Security event detected", fields)
}

// LogDatabaseOperation logs database operations
func LogDatabaseOperation(operation, table string, recordID interface{}, duration time.Duration, err error) {
	fields := map[string]interface{}{
		"db_operation": operation,
		"table":        table,
		"duration_ms":  duration.Milliseconds(),
	}
	
	if recordID != nil {
		fields["record_id"] = recordID
	}
	
	if err != nil {
		fields["error"] = err.Error()
		Error("Database operation failed", err, fields)
	} else {
		Debug("Database operation completed", fields)
	}
}

// LogBusinessEvent logs business-related events
func LogBusinessEvent(event, entityType string, entityID interface{}, details map[string]interface{}) {
	fields := map[string]interface{}{
		"business_event": event,
		"entity_type":    entityType,
	}
	
	if entityID != nil {
		fields["entity_id"] = entityID
	}
	
	for k, v := range details {
		fields[k] = v
	}
	
	Info("Business event logged", fields)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
package auth

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GetJWTSecret returns the JWT secret from environment
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "blog-service-default-secret-change-this-in-production"
	}
	return []byte(secret)
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	// Check if header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	// Extract token part
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

// ValidateAccessToken validates JWT token and returns claims
func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	// Parse token with claims
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Additional validation
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateAccessToken generates a new JWT access token
func GenerateAccessToken(userID uint, email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24 hours default

	// Parse duration from environment if available
	if durationStr := os.Getenv("JWT_ACCESS_TOKEN_DURATION"); durationStr != "" {
		if duration, err := time.ParseDuration(durationStr); err == nil {
			expirationTime = time.Now().Add(duration)
		}
	}

	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    getEnv("JWT_ISSUER", "mejona-blog-service"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// HasPermission checks if a role has a specific permission
func HasPermission(role, permission string) bool {
	// Simple permission system - can be enhanced later
	permissions := map[string][]string{
		"admin": {
			"blog:create", "blog:read", "blog:update", "blog:delete",
			"blog:publish", "blog:unpublish", "blog:moderate",
			"user:create", "user:read", "user:update", "user:delete",
			"admin:all",
		},
		"manager": {
			"blog:create", "blog:read", "blog:update", "blog:delete",
			"blog:publish", "blog:unpublish",
			"user:read", "user:update",
		},
		"editor": {
			"blog:create", "blog:read", "blog:update",
			"blog:publish", "blog:unpublish",
		},
		"author": {
			"blog:create", "blog:read", "blog:update",
		},
		"user": {
			"blog:read",
		},
	}

	rolePermissions, exists := permissions[role]
	if !exists {
		return false
	}

	// Check for wildcard permission
	for _, perm := range rolePermissions {
		if perm == "admin:all" || perm == permission {
			return true
		}
	}

	return false
}

// GetRoleHierarchy returns role hierarchy level (higher number = more privileges)
func GetRoleHierarchy(role string) int {
	hierarchy := map[string]int{
		"admin":   5,
		"manager": 4,
		"editor":  3,
		"author":  2,
		"user":    1,
		"guest":   0,
	}

	level, exists := hierarchy[role]
	if !exists {
		return 0 // Default to guest level
	}

	return level
}

// HasRoleOrAbove checks if user has the specified role or higher
func HasRoleOrAbove(userRole, requiredRole string) bool {
	return GetRoleHierarchy(userRole) >= GetRoleHierarchy(requiredRole)
}

// ValidateRole checks if a role is valid
func ValidateRole(role string) bool {
	validRoles := []string{"admin", "manager", "editor", "author", "user"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

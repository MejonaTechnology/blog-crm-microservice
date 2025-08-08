# Blog CRM Management Microservice

A professional Go microservice for blog content management and analytics, part of Mejona Technology's admin dashboard system.

## Service Information

- **Port**: 8082
- **Framework**: Gin (Go)
- **Database**: MySQL (mejona_unified)
- **Architecture**: Microservice with Clean Architecture patterns

## Project Structure

```
blog-service/
├── cmd/server/           # Application entry point
│   └── main.go          # Main server file with health endpoints
├── internal/            # Internal application code
│   ├── handlers/        # HTTP handlers
│   │   └── health.go   # Health check endpoints
│   ├── middleware/      # HTTP middleware
│   │   ├── auth.go     # JWT authentication
│   │   ├── cors.go     # CORS handling
│   │   └── logging.go  # Request logging
│   ├── models/          # Data models (to be implemented)
│   └── services/        # Business logic (to be implemented)
├── pkg/                 # Shared packages
│   ├── auth/           # Authentication utilities
│   │   └── auth.go     # JWT token handling
│   ├── database/       # Database connection
│   │   └── database.go # GORM MySQL connection
│   └── logger/         # Structured logging
│       └── logger.go   # Logrus-based logger
├── .env.example        # Environment configuration template
├── .env                # Local environment variables
├── go.mod              # Go module dependencies
└── README.md           # This file
```

## Available Endpoints

### Health Check Endpoints
- `GET /health` - Basic health check
- `GET /health/deep` - Comprehensive health check with database
- `GET /status` - Quick status information
- `GET /ready` - Readiness probe for Kubernetes
- `GET /alive` - Liveness probe for Kubernetes
- `GET /metrics` - System and performance metrics

### API Endpoints
- `GET /api/v1/test` - Test endpoint for service verification

### Documentation
- `GET /swagger/index.html` - Swagger API documentation (if enabled)

## Configuration

### Environment Variables

Key configuration variables (see `.env.example` for complete list):

```bash
# Database Configuration
DB_HOST=localhost
DB_NAME=mejona_unified
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_PORT=3306

# Server Configuration
PORT=8082
GIN_MODE=debug
APP_DEBUG=true

# JWT Configuration
JWT_SECRET=your-jwt-secret-key
JWT_ACCESS_TOKEN_DURATION=24h

# Logging Configuration
LOG_LEVEL=info
LOG_FILE_ENABLED=true
LOG_FILE_PATH=./logs/blog-service.log
```

## Development

### Prerequisites

- Go 1.23.0 or higher
- MySQL 8.0 or higher
- Access to `mejona_unified` database

### Building

```bash
# Install dependencies
go mod tidy

# Build the service
go build -o blog-service.exe ./cmd/server/main.go
```

### Running

```bash
# Start the service
./blog-service.exe

# Or use the test script
test-service.bat
```

### Testing

The service includes a test script that verifies:
- Successful compilation
- Service startup (will show database connection error if MySQL not available)
- Basic service architecture validation

## Features Implemented

✅ **Core Architecture**
- Clean Architecture with separation of concerns
- Gin HTTP framework with middleware pipeline
- GORM database integration with connection pooling
- Structured logging with file rotation

✅ **Security**
- JWT-based authentication with configurable expiration
- CORS middleware with environment-based origin control
- Role-based access control (admin, manager, editor, author, user)
- Security event logging

✅ **Monitoring & Health**
- Multiple health check endpoints for different use cases
- Comprehensive system metrics (memory, database, performance)
- Request/response logging with correlation IDs
- Panic recovery with security event logging

✅ **Configuration**
- Environment-based configuration
- Database connection pooling with configurable limits
- Feature flags for enabling/disabling functionality
- Development vs production environment support

## Integration Points

This service is designed to integrate with:
- **Contact Service** (port 8081) - CRM integration
- **Admin Backend** (port 8080) - Main dashboard API
- **React Frontend** (port 5173) - Admin dashboard UI

## Next Steps

The foundation is complete and ready for feature implementation:

1. **Blog Models** - Create blog, category, tag, comment models
2. **Blog Handlers** - CRUD operations for blog management
3. **Analytics Services** - Blog performance and engagement tracking
4. **SEO Services** - SEO analysis and optimization features
5. **Publishing Workflow** - Draft → Review → Publish pipeline
6. **Comment Moderation** - Comment approval and spam detection

## Status

✅ **FOUNDATION COMPLETE** - Service builds and runs successfully
- All core middleware implemented
- Health endpoints functional
- Database integration ready
- Authentication system ready
- Logging and monitoring configured

The blog-service is now ready for feature development and can be deployed alongside other microservices in the Mejona Technology admin dashboard ecosystem.# Trigger deployment with configured secrets

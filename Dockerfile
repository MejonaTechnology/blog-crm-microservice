# Multi-stage Docker build for Blog CRM Microservice
# Stage 1: Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o blog-service \
    ./cmd/server/main.go

# Verify the binary was created
RUN ls -la blog-service

# Stage 2: Final stage
FROM alpine:latest

# Install ca-certificates, tzdata, and curl for health checks
RUN apk --no-cache add ca-certificates tzdata curl wget

# Create non-root user
RUN adduser -D -s /bin/sh -u 1002 bloguser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/blog-service .

# Copy configuration and migration files
COPY --from=builder /app/migrations ./migrations/
COPY --from=builder /app/docs ./docs/
COPY --from=builder /app/config ./config/

# Create necessary directories
RUN mkdir -p logs uploads temp cache

# Set ownership to bloguser
RUN chown -R bloguser:bloguser /app

# Make binary executable
RUN chmod +x /app/blog-service

# Switch to non-root user
USER bloguser

# Expose port 8082
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl --fail --silent --show-error http://localhost:8082/health || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8082
ENV TZ=Asia/Kolkata

# Labels for better maintainability
LABEL maintainer="Mejona Technology <support@mejona.com>"
LABEL version="1.0.0"
LABEL description="Blog CRM Management Microservice for Mejona Technology"
LABEL service="blog-service"
LABEL port="8082"

# Run the application
CMD ["./blog-service"]
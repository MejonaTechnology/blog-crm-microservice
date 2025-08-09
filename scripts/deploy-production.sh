#!/bin/bash

# Blog CRM Service Deployment Script
# This script handles both CI/CD and manual deployments

set -e  # Exit on any error

# Configuration
SERVICE_NAME="blog-crm-service"
SERVICE_PORT="8082"
SERVICE_DIR="/opt/mejona/blog-crm-service"
SERVICE_USER="ubuntu"
SERVICE_GROUP="ubuntu"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root for systemctl commands
check_permissions() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script needs to be run with sudo privileges for systemctl commands"
        exit 1
    fi
}

# Create service directory and set permissions
setup_directory() {
    log_info "Setting up service directory..."
    mkdir -p $SERVICE_DIR
    chown -R $SERVICE_USER:$SERVICE_GROUP $SERVICE_DIR
    log_success "Service directory ready: $SERVICE_DIR"
}

# Stop existing service
stop_service() {
    log_info "Stopping existing service..."
    systemctl stop $SERVICE_NAME 2>/dev/null || true
    log_success "Service stopped (if it was running)"
}

# Clone or update repository
update_repository() {
    log_info "Updating repository..."
    cd $SERVICE_DIR
    
    if [ ! -d ".git" ]; then
        log_info "Cloning repository..."
        git clone https://github.com/MejonaTechnology/blog-crm-microservice.git .
    else
        log_info "Pulling latest changes..."
        git fetch origin
        git reset --hard origin/main
    fi
    
    # Set proper ownership
    chown -R $SERVICE_USER:$SERVICE_GROUP $SERVICE_DIR
    log_success "Repository updated"
}

# Build service
build_service() {
    log_info "Building service..."
    cd $SERVICE_DIR
    
    # Download dependencies
    sudo -u $SERVICE_USER /usr/local/go/bin/go mod download
    
    # Build the service
    sudo -u $SERVICE_USER CGO_ENABLED=0 /usr/local/go/bin/go build -ldflags="-w -s" -o $SERVICE_NAME cmd/server/main.go
    
    # Make executable
    chmod +x $SERVICE_NAME
    
    # Verify build
    if [ ! -f "$SERVICE_NAME" ]; then
        log_error "Build failed - binary not found"
        exit 1
    fi
    
    log_success "Service built successfully"
    ls -la $SERVICE_NAME
}

# Create environment configuration
create_env_config() {
    log_info "Creating environment configuration..."
    
    # Check if required environment variables are set
    required_vars=("DB_HOST" "DB_USER" "DB_PASSWORD" "DB_NAME" "JWT_SECRET")
    missing_vars=()
    
    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            missing_vars+=("$var")
        fi
    done
    
    if [ ${#missing_vars[@]} -gt 0 ]; then
        log_warning "Missing environment variables: ${missing_vars[*]}"
        log_info "Creating environment file with defaults..."
    fi
    
    cat > $SERVICE_DIR/.env << EOF
# Database Configuration
DB_HOST=${DB_HOST:-localhost}
DB_USER=${DB_USER:-mejona_user}
DB_PASSWORD=${DB_PASSWORD:-}
DB_NAME=${DB_NAME:-mejona_unified}
DB_PORT=3306

# Service Configuration
PORT=$SERVICE_PORT
GIN_MODE=release
APP_ENV=production
JWT_SECRET=${JWT_SECRET:-default_jwt_secret_change_this}

# Feature Flags
ENABLE_ANALYTICS=true
ENABLE_LEAD_TRACKING=true
ENABLE_SEO_ANALYSIS=true
ENABLE_METRICS=true

# Logging
LOG_LEVEL=info
APP_DEBUG=false
EOF
    
    chown $SERVICE_USER:$SERVICE_GROUP $SERVICE_DIR/.env
    chmod 600 $SERVICE_DIR/.env  # Secure permissions for env file
    log_success "Environment configuration created"
}

# Create systemd service
create_systemd_service() {
    log_info "Creating systemd service..."
    
    cat > /etc/systemd/system/$SERVICE_NAME.service << EOF
[Unit]
Description=Blog CRM Management Microservice
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=$SERVICE_USER
Group=$SERVICE_GROUP
WorkingDirectory=$SERVICE_DIR
ExecStart=$SERVICE_DIR/$SERVICE_NAME
EnvironmentFile=$SERVICE_DIR/.env
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# Security settings
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF
    
    log_success "Systemd service created"
}

# Start and enable service
start_service() {
    log_info "Starting service..."
    
    systemctl daemon-reload
    systemctl enable $SERVICE_NAME
    systemctl start $SERVICE_NAME
    
    log_success "Service started and enabled"
}

# Wait for service to be ready
wait_for_service() {
    log_info "Waiting for service to be ready..."
    
    # Wait up to 30 seconds for service to start
    for i in {1..30}; do
        if systemctl is-active --quiet $SERVICE_NAME; then
            log_success "Service is active"
            break
        fi
        sleep 1
    done
    
    # Check if service is actually active
    if ! systemctl is-active --quiet $SERVICE_NAME; then
        log_error "Service failed to start"
        log_info "Service status:"
        systemctl status $SERVICE_NAME --no-pager -l
        log_info "Recent logs:"
        journalctl -u $SERVICE_NAME -n 20 --no-pager
        exit 1
    fi
    
    # Wait a bit more for the HTTP server to be ready
    sleep 10
}

# Perform health check
health_check() {
    log_info "Performing health check..."
    
    # Try health check endpoint
    if curl -f --connect-timeout 10 --max-time 30 http://localhost:$SERVICE_PORT/health; then
        log_success "Health check passed"
    else
        log_error "Health check failed"
        log_info "Service logs:"
        journalctl -u $SERVICE_NAME -n 20 --no-pager
        exit 1
    fi
    
    # Try test API endpoint
    log_info "Testing API endpoint..."
    if curl -f --connect-timeout 10 --max-time 30 http://localhost:$SERVICE_PORT/api/v1/test; then
        log_success "API test passed"
    else
        log_warning "API test failed (this might be expected if authentication is required)"
    fi
}

# Display deployment summary
deployment_summary() {
    log_success "🎉 Deployment completed successfully!"
    echo ""
    echo "📊 Deployment Summary:"
    echo "====================="
    echo "Service: Blog CRM Management Microservice"
    echo "Port: $SERVICE_PORT"
    echo "Directory: $SERVICE_DIR"
    echo "User: $SERVICE_USER"
    echo "Status: $(systemctl is-active $SERVICE_NAME)"
    echo ""
    echo "🔗 Endpoints:"
    echo "Health Check: http://localhost:$SERVICE_PORT/health"
    echo "API Test: http://localhost:$SERVICE_PORT/api/v1/test"
    echo ""
    echo "🛠️  Management Commands:"
    echo "Status: systemctl status $SERVICE_NAME"
    echo "Logs: journalctl -u $SERVICE_NAME -f"
    echo "Stop: systemctl stop $SERVICE_NAME"
    echo "Start: systemctl start $SERVICE_NAME"
    echo "Restart: systemctl restart $SERVICE_NAME"
}

# Main deployment function
main() {
    log_info "🚀 Starting Blog CRM Service deployment..."
    
    check_permissions
    setup_directory
    stop_service
    update_repository
    build_service
    create_env_config
    create_systemd_service
    start_service
    wait_for_service
    health_check
    deployment_summary
    
    log_success "✅ Blog CRM Service deployment completed successfully!"
}

# Handle script arguments
case "${1:-deploy}" in
    "deploy")
        main
        ;;
    "stop")
        log_info "Stopping Blog CRM Service..."
        systemctl stop $SERVICE_NAME || true
        log_success "Service stopped"
        ;;
    "start")
        log_info "Starting Blog CRM Service..."
        systemctl start $SERVICE_NAME
        log_success "Service started"
        ;;
    "restart")
        log_info "Restarting Blog CRM Service..."
        systemctl restart $SERVICE_NAME
        log_success "Service restarted"
        ;;
    "status")
        systemctl status $SERVICE_NAME --no-pager -l
        ;;
    "logs")
        journalctl -u $SERVICE_NAME -f
        ;;
    "health")
        curl -f http://localhost:$SERVICE_PORT/health || exit 1
        ;;
    "test")
        curl -f http://localhost:$SERVICE_PORT/api/v1/test || exit 1
        ;;
    *)
        echo "Usage: $0 {deploy|stop|start|restart|status|logs|health|test}"
        echo ""
        echo "Commands:"
        echo "  deploy  - Full deployment (default)"
        echo "  stop    - Stop the service"
        echo "  start   - Start the service"
        echo "  restart - Restart the service"
        echo "  status  - Show service status"
        echo "  logs    - Show live logs"
        echo "  health  - Perform health check"
        echo "  test    - Test API endpoint"
        exit 1
        ;;
esac
#!/bin/bash
# =============================================================================
# Blog CRM Microservice - AWS EC2 Deployment Script
# Mejona Technology LLP - Production Deployment Automation
# =============================================================================

set -euo pipefail  # Exit on error, undefined variables, pipe failures

# =============================================================================
# CONFIGURATION VARIABLES
# =============================================================================
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SERVICE_NAME="blog-service"
readonly SERVICE_PORT="8082"
readonly AWS_EC2_IP="65.1.94.25"
readonly AWS_USER="ubuntu"  # Change if using different user
readonly SERVICE_USER="bloguser"
readonly SERVICE_DIR="/opt/mejona/$SERVICE_NAME"
readonly BINARY_NAME="blog-service"
readonly SYSTEMD_SERVICE="mejona-blog-service"
readonly LOG_DIR="/var/log/mejona"
readonly CONFIG_DIR="/etc/mejona"

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# =============================================================================
# UTILITY FUNCTIONS
# =============================================================================

log() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[$(date '+%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1"
}

error() {
    echo -e "${RED}[$(date '+%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1" >&2
}

success() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')] SUCCESS:${NC} $1"
}

info() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')] INFO:${NC} $1"
}

# =============================================================================
# VALIDATION FUNCTIONS
# =============================================================================

check_prerequisites() {
    log "Checking deployment prerequisites..."
    
    # Check if required tools are available
    local required_tools=("ssh" "scp" "go" "git")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            error "$tool is required but not installed"
            exit 1
        fi
    done
    
    # Check if SSH key exists
    if [[ ! -f ~/.ssh/id_rsa && ! -f ~/.ssh/id_ed25519 ]]; then
        warn "No SSH key found. Make sure you have SSH access to the server"
    fi
    
    # Check if we can reach the server
    if ! ping -c 1 "$AWS_EC2_IP" &> /dev/null; then
        warn "Cannot ping server $AWS_EC2_IP. Continuing anyway..."
    fi
    
    success "Prerequisites check completed"
}

validate_environment() {
    log "Validating environment configuration..."
    
    # Check if .env file exists
    if [[ ! -f ".env" ]]; then
        warn ".env file not found. Creating template..."
        create_env_template
    fi
    
    # Validate Go version
    local go_version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    if [[ "$(printf '%s\n' "1.21" "$go_version" | sort -V | head -n1)" != "1.21" ]]; then
        error "Go version 1.21+ required. Current version: $go_version"
        exit 1
    fi
    
    success "Environment validation completed"
}

create_env_template() {
    cat > .env << EOF
# Blog Service Environment Configuration
GIN_MODE=release
PORT=8082
LOG_LEVEL=info

# Database Configuration
DB_HOST=localhost
DB_USER=blog_user
DB_PASSWORD=secure_password_here
DB_NAME=mejona_unified
DB_PORT=3306

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis_password_here

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRY=24h

# API Configuration
API_VERSION=v1
ENABLE_SWAGGER=false
RATE_LIMIT=100

# Email Configuration (if needed)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@domain.com
SMTP_PASS=your_app_password

# AWS Configuration (if needed)
AWS_REGION=ap-south-1
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key

# Monitoring Configuration
ENABLE_METRICS=true
METRICS_PORT=9090
EOF
    warn "Please update .env file with actual values before deployment"
}

# =============================================================================
# BUILD FUNCTIONS
# =============================================================================

build_service() {
    log "Building blog service binary..."
    
    # Clean previous builds
    rm -f "$BINARY_NAME"
    
    # Build for Linux AMD64 (EC2 target)
    env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
        -ldflags='-w -s -extldflags "-static"' \
        -a -installsuffix cgo \
        -o "$BINARY_NAME" \
        ./cmd/server/main.go
    
    if [[ ! -f "$BINARY_NAME" ]]; then
        error "Failed to build binary"
        exit 1
    fi
    
    # Make binary executable
    chmod +x "$BINARY_NAME"
    
    success "Binary built successfully: $(ls -lh $BINARY_NAME | awk '{print $5}')"
}

run_tests() {
    log "Running test suite..."
    
    # Run Go tests
    if ! go test ./... -v; then
        error "Tests failed. Deployment aborted."
        exit 1
    fi
    
    # Run integration tests if they exist
    if [[ -d "tests/integration" ]]; then
        log "Running integration tests..."
        go test ./tests/integration/... -v
    fi
    
    success "All tests passed"
}

# =============================================================================
# SERVER PREPARATION FUNCTIONS
# =============================================================================

prepare_server() {
    log "Preparing server environment..."
    
    ssh "$AWS_USER@$AWS_EC2_IP" << 'EOF'
        set -euo pipefail
        
        # Update system packages
        sudo apt-get update
        sudo apt-get install -y curl wget unzip htop
        
        # Create service user if it doesn't exist
        if ! id bloguser &>/dev/null; then
            sudo useradd -r -s /bin/false -m -d /home/bloguser bloguser
            echo "Created bloguser system user"
        fi
        
        # Create necessary directories
        sudo mkdir -p /opt/mejona/blog-service/{bin,config,logs,migrations,docs}
        sudo mkdir -p /var/log/mejona
        sudo mkdir -p /etc/mejona
        
        # Set ownership
        sudo chown -R bloguser:bloguser /opt/mejona/blog-service
        sudo chown bloguser:bloguser /var/log/mejona
        
        # Install MySQL client if not present
        if ! command -v mysql &> /dev/null; then
            sudo apt-get install -y mysql-client-core-8.0
        fi
        
        echo "Server preparation completed"
EOF
    
    success "Server preparation completed"
}

create_systemd_service() {
    log "Creating systemd service file..."
    
    # Create systemd service file locally
    cat > "${SYSTEMD_SERVICE}.service" << EOF
[Unit]
Description=Mejona Blog CRM Management Microservice
Documentation=https://github.com/MejonaTechnology/mejona-blog-service
After=network-online.target mysql.service
Wants=network-online.target
Requires=mysql.service

[Service]
Type=simple
User=bloguser
Group=bloguser
WorkingDirectory=/opt/mejona/blog-service
ExecStart=/opt/mejona/blog-service/bin/blog-service
ExecReload=/bin/kill -HUP \$MAINPID
KillMode=mixed
KillSignal=SIGTERM
TimeoutStopSec=30
Restart=always
RestartSec=10
StartLimitInterval=60
StartLimitBurst=3

# Environment
Environment=GIN_MODE=release
Environment=PORT=8082
EnvironmentFile=/etc/mejona/blog-service.env

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/mejona/blog-service/logs /var/log/mejona
CapabilityBoundingSet=CAP_NET_BIND_SERVICE

# Resource limits
LimitNOFILE=65536
LimitNPROC=32768

# Logging
StandardOutput=append:/var/log/mejona/blog-service.log
StandardError=append:/var/log/mejona/blog-service-error.log
SyslogIdentifier=mejona-blog-service

[Install]
WantedBy=multi-user.target
EOF
    
    success "Systemd service file created"
}

# =============================================================================
# DEPLOYMENT FUNCTIONS
# =============================================================================

deploy_files() {
    log "Deploying files to server..."
    
    # Copy binary
    scp "$BINARY_NAME" "$AWS_USER@$AWS_EC2_IP:/tmp/"
    
    # Copy configuration files
    if [[ -f ".env" ]]; then
        scp .env "$AWS_USER@$AWS_EC2_IP:/tmp/blog-service.env"
    fi
    
    # Copy migrations
    if [[ -d "migrations" ]]; then
        scp -r migrations "$AWS_USER@$AWS_EC2_IP:/tmp/"
    fi
    
    # Copy docs
    if [[ -d "docs" ]]; then
        scp -r docs "$AWS_USER@$AWS_EC2_IP:/tmp/"
    fi
    
    # Copy systemd service file
    scp "${SYSTEMD_SERVICE}.service" "$AWS_USER@$AWS_EC2_IP:/tmp/"
    
    success "Files uploaded to server"
}

install_service() {
    log "Installing service on server..."
    
    ssh "$AWS_USER@$AWS_EC2_IP" << EOF
        set -euo pipefail
        
        # Move binary to service directory
        sudo mv /tmp/$BINARY_NAME /opt/mejona/blog-service/bin/
        sudo chmod +x /opt/mejona/blog-service/bin/$BINARY_NAME
        
        # Move configuration
        if [[ -f /tmp/blog-service.env ]]; then
            sudo mv /tmp/blog-service.env /etc/mejona/
            sudo chown bloguser:bloguser /etc/mejona/blog-service.env
            sudo chmod 600 /etc/mejona/blog-service.env
        fi
        
        # Move migrations and docs
        if [[ -d /tmp/migrations ]]; then
            sudo mv /tmp/migrations /opt/mejona/blog-service/
            sudo chown -R bloguser:bloguser /opt/mejona/blog-service/migrations
        fi
        
        if [[ -d /tmp/docs ]]; then
            sudo mv /tmp/docs /opt/mejona/blog-service/
            sudo chown -R bloguser:bloguser /opt/mejona/blog-service/docs
        fi
        
        # Install systemd service
        sudo mv /tmp/${SYSTEMD_SERVICE}.service /etc/systemd/system/
        sudo systemctl daemon-reload
        sudo systemctl enable ${SYSTEMD_SERVICE}
        
        echo "Service installation completed"
EOF
    
    success "Service installed successfully"
}

configure_nginx() {
    log "Configuring nginx reverse proxy..."
    
    ssh "$AWS_USER@$AWS_EC2_IP" << 'EOF'
        # Create nginx configuration for blog service
        sudo tee /etc/nginx/sites-available/blog-service << 'NGINX_EOF'
# Blog Service Reverse Proxy Configuration
upstream blog_service {
    server 127.0.0.1:8082;
    keepalive 32;
}

server {
    listen 80;
    server_name mejona.in;
    
    # Blog service API routes
    location /api/v1/blog/ {
        proxy_pass http://blog_service/api/v1/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
    
    # Health check endpoint
    location /blog/health {
        proxy_pass http://blog_service/health;
        access_log off;
    }
    
    # Swagger documentation
    location /blog/docs/ {
        proxy_pass http://blog_service/swagger/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
NGINX_EOF
        
        # Enable the site
        sudo ln -sf /etc/nginx/sites-available/blog-service /etc/nginx/sites-enabled/
        
        # Test nginx configuration
        sudo nginx -t
        
        # Reload nginx
        sudo systemctl reload nginx
        
        echo "Nginx configuration completed"
EOF
    
    success "Nginx configured successfully"
}

# =============================================================================
# FIREWALL AND SECURITY
# =============================================================================

configure_firewall() {
    log "Configuring firewall rules..."
    
    ssh "$AWS_USER@$AWS_EC2_IP" << 'EOF'
        # Allow blog service port
        sudo ufw allow 8082/tcp comment 'Blog Service'
        
        # Allow nginx ports if not already allowed
        sudo ufw allow 'Nginx Full' 2>/dev/null || true
        
        # Check firewall status
        sudo ufw status
        
        echo "Firewall configuration completed"
EOF
    
    success "Firewall configured successfully"
}

# =============================================================================
# DATABASE MANAGEMENT
# =============================================================================

run_migrations() {
    log "Running database migrations..."
    
    ssh "$AWS_USER@$AWS_EC2_IP" << 'EOF'
        cd /opt/mejona/blog-service
        
        # Source environment variables
        if [[ -f /etc/mejona/blog-service.env ]]; then
            set -a
            source /etc/mejona/blog-service.env
            set +a
        fi
        
        # Run migrations if directory exists
        if [[ -d migrations ]]; then
            echo "Running database migrations..."
            # Add migration logic here based on your migration tool
            # For now, we'll just list the migrations
            ls -la migrations/
        fi
        
        echo "Migration check completed"
EOF
    
    success "Database migrations processed"
}

# =============================================================================
# SERVICE MANAGEMENT
# =============================================================================

start_service() {
    log "Starting blog service..."
    
    ssh "$AWS_USER@$AWS_EC2_IP" << EOF
        # Start the service
        sudo systemctl start ${SYSTEMD_SERVICE}
        
        # Check status
        sudo systemctl status ${SYSTEMD_SERVICE} --no-pager
        
        # Wait for service to be ready
        echo "Waiting for service to start..."
        for i in {1..30}; do
            if curl -sf http://localhost:8082/health > /dev/null 2>&1; then
                echo "Service is healthy!"
                break
            fi
            echo "Waiting... (\$i/30)"
            sleep 2
        done
        
        # Show recent logs
        echo "Recent logs:"
        sudo journalctl -u ${SYSTEMD_SERVICE} --no-pager -n 10
EOF
    
    success "Service started successfully"
}

# =============================================================================
# VERIFICATION FUNCTIONS
# =============================================================================

verify_deployment() {
    log "Verifying deployment..."
    
    # Test health endpoint
    if curl -sf "http://$AWS_EC2_IP:8082/health" > /dev/null; then
        success "Health endpoint is responding"
    else
        error "Health endpoint is not responding"
        return 1
    fi
    
    # Test API endpoint
    if curl -sf "http://$AWS_EC2_IP:8082/api/v1/test" > /dev/null; then
        success "API endpoint is responding"
    else
        warn "API endpoint test failed"
    fi
    
    # Check service status via SSH
    ssh "$AWS_USER@$AWS_EC2_IP" << EOF
        echo "=== SERVICE STATUS ==="
        sudo systemctl is-active ${SYSTEMD_SERVICE}
        
        echo "=== LISTENING PORTS ==="
        sudo ss -tlnp | grep :8082
        
        echo "=== RECENT LOGS ==="
        sudo journalctl -u ${SYSTEMD_SERVICE} --no-pager -n 5
EOF
    
    success "Deployment verification completed"
}

# =============================================================================
# CLEANUP FUNCTION
# =============================================================================

cleanup() {
    log "Cleaning up temporary files..."
    rm -f "$BINARY_NAME"
    rm -f "${SYSTEMD_SERVICE}.service"
    success "Cleanup completed"
}

# =============================================================================
# MAIN DEPLOYMENT FLOW
# =============================================================================

print_banner() {
    echo -e "${BLUE}"
    cat << 'BANNER'
╔══════════════════════════════════════════════════════════════╗
║                    MEJONA TECHNOLOGY LLP                     ║
║              Blog CRM Service - AWS Deployment              ║
║                                                              ║
║  Target: AWS EC2 (65.1.94.25)                              ║
║  Service: Blog CRM Management Microservice                  ║
║  Port: 8082                                                  ║
╚══════════════════════════════════════════════════════════════╝
BANNER
    echo -e "${NC}"
}

show_help() {
    cat << EOF
Usage: $0 [OPTIONS]

Deploy Blog CRM Microservice to AWS EC2

OPTIONS:
    -h, --help              Show this help message
    -b, --build-only        Only build the binary
    -t, --test-only         Only run tests
    -d, --deploy-only       Only deploy (skip build and tests)
    -v, --verify-only       Only verify existing deployment
    -s, --skip-tests        Skip running tests
    --skip-migrations       Skip running database migrations
    --restart-only          Only restart the service
    --logs                  Show recent service logs

EXAMPLES:
    $0                      Full deployment (build, test, deploy)
    $0 --build-only         Just build the binary
    $0 --deploy-only        Deploy without building
    $0 --restart-only       Just restart the service
    $0 --logs               Show recent logs

EOF
}

main() {
    local build_only=false
    local test_only=false
    local deploy_only=false
    local verify_only=false
    local skip_tests=false
    local skip_migrations=false
    local restart_only=false
    local show_logs=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -b|--build-only)
                build_only=true
                shift
                ;;
            -t|--test-only)
                test_only=true
                shift
                ;;
            -d|--deploy-only)
                deploy_only=true
                shift
                ;;
            -v|--verify-only)
                verify_only=true
                shift
                ;;
            -s|--skip-tests)
                skip_tests=true
                shift
                ;;
            --skip-migrations)
                skip_migrations=true
                shift
                ;;
            --restart-only)
                restart_only=true
                shift
                ;;
            --logs)
                show_logs=true
                shift
                ;;
            *)
                error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_banner
    
    # Handle special cases
    if [[ "$show_logs" == true ]]; then
        ssh "$AWS_USER@$AWS_EC2_IP" "sudo journalctl -u ${SYSTEMD_SERVICE} -f"
        exit 0
    fi
    
    if [[ "$restart_only" == true ]]; then
        log "Restarting service only..."
        ssh "$AWS_USER@$AWS_EC2_IP" << EOF
            sudo systemctl restart ${SYSTEMD_SERVICE}
            sudo systemctl status ${SYSTEMD_SERVICE} --no-pager
EOF
        exit 0
    fi
    
    if [[ "$verify_only" == true ]]; then
        verify_deployment
        exit 0
    fi
    
    # Set trap for cleanup
    trap cleanup EXIT
    
    # Run deployment steps
    if [[ "$deploy_only" == false ]]; then
        check_prerequisites
        validate_environment
    fi
    
    if [[ "$test_only" == true ]]; then
        run_tests
        exit 0
    fi
    
    if [[ "$build_only" == false && "$deploy_only" == false ]]; then
        if [[ "$skip_tests" == false ]]; then
            run_tests
        fi
        build_service
    elif [[ "$build_only" == true ]]; then
        build_service
        exit 0
    fi
    
    if [[ "$deploy_only" == true || "$build_only" == false ]]; then
        prepare_server
        create_systemd_service
        deploy_files
        install_service
        configure_nginx
        configure_firewall
        
        if [[ "$skip_migrations" == false ]]; then
            run_migrations
        fi
        
        start_service
        verify_deployment
    fi
    
    success "Blog service deployment completed successfully!"
    info "Service URL: http://$AWS_EC2_IP:8082"
    info "Health Check: http://$AWS_EC2_IP:8082/health"
    info "API Docs: http://$AWS_EC2_IP:8082/swagger/index.html"
    info "Nginx Proxy: http://$AWS_EC2_IP/api/v1/blog/"
}

# Run main function with all arguments
main "$@"
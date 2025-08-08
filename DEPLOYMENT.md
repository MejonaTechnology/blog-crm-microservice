# Blog CRM Microservice - Deployment Guide

![Mejona Technology](https://img.shields.io/badge/Mejona-Technology-blue)
![Version](https://img.shields.io/badge/version-1.0.0-green)
![Go](https://img.shields.io/badge/go-1.23+-blue)
![Port](https://img.shields.io/badge/port-8082-orange)

## Overview

This guide provides comprehensive instructions for deploying the Blog CRM Management Microservice in production environments. The service is designed to run on **AWS EC2 (65.1.94.25:8082)** alongside existing Mejona services.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    AWS EC2 (65.1.94.25)                    │
├─────────────────────────────────────────────────────────────┤
│  Nginx (80/443) → Blog Service (8082)                      │
│  ├── /api/v1/blog/* → http://localhost:8082/api/v1/*       │
│  ├── /blog/health   → http://localhost:8082/health         │
│  └── /blog/docs     → http://localhost:8082/swagger        │
├─────────────────────────────────────────────────────────────┤
│  Services:                                                  │
│  ├── mejona-blog-service (systemd)                        │
│  ├── contact-service (8081)                               │
│  └── admin-backend (8080)                                 │
├─────────────────────────────────────────────────────────────┤
│  Database: mejona_unified (MySQL 8.0)                     │
│  Cache: Redis (optional)                                  │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Deployment

### Prerequisites

- **Server**: AWS EC2 with Ubuntu 20.04+
- **Go**: Version 1.23+
- **Database**: MySQL 8.0 (mejona_unified database)
- **Web Server**: Nginx (for reverse proxy)
- **User Access**: SSH access to ubuntu@65.1.94.25

### One-Click Deployment

```bash
# Clone the repository
git clone https://github.com/MejonaTechnology/blog-service.git
cd blog-service

# Make deployment script executable
chmod +x scripts/deploy-blog-service.sh

# Run full deployment (build, test, deploy)
./scripts/deploy-blog-service.sh

# Or deploy only (skip build and tests)
./scripts/deploy-blog-service.sh --deploy-only
```

## 📋 Manual Deployment Steps

### 1. Environment Setup

Create environment configuration:

```bash
# Create .env file
cat > .env << EOF
# Service Configuration
GIN_MODE=release
PORT=8082
LOG_LEVEL=info

# Database Configuration
DB_HOST=localhost
DB_USER=blog_user
DB_PASSWORD=your_secure_password
DB_NAME=mejona_unified
DB_PORT=3306

# Redis Configuration (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_256_chars
JWT_EXPIRY=24h

# API Configuration
API_VERSION=v1
ENABLE_SWAGGER=true
RATE_LIMIT=100

# Monitoring
ENABLE_METRICS=true
METRICS_PORT=9090
EOF
```

### 2. Build Service

```bash
# Install dependencies
go mod download

# Run tests (optional but recommended)
go test -v ./...

# Build for Linux (if building on different OS)
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o blog-service \
    ./cmd/server/main.go
```

### 3. Server Preparation

```bash
# SSH to server
ssh ubuntu@65.1.94.25

# Create service user
sudo useradd -r -s /bin/false -m -d /home/bloguser bloguser

# Create directories
sudo mkdir -p /opt/mejona/blog-service/{bin,config,logs,migrations,docs}
sudo mkdir -p /var/log/mejona
sudo mkdir -p /etc/mejona

# Set permissions
sudo chown -R bloguser:bloguser /opt/mejona/blog-service
sudo chown bloguser:bloguser /var/log/mejona
```

### 4. Deploy Files

```bash
# Copy binary
scp blog-service ubuntu@65.1.94.25:/tmp/
ssh ubuntu@65.1.94.25 "sudo mv /tmp/blog-service /opt/mejona/blog-service/bin/ && sudo chmod +x /opt/mejona/blog-service/bin/blog-service"

# Copy configuration
scp .env ubuntu@65.1.94.25:/tmp/blog-service.env
ssh ubuntu@65.1.94.25 "sudo mv /tmp/blog-service.env /etc/mejona/ && sudo chown bloguser:bloguser /etc/mejona/blog-service.env"

# Copy systemd service
scp scripts/mejona-blog-service.service ubuntu@65.1.94.25:/tmp/
ssh ubuntu@65.1.94.25 "sudo mv /tmp/mejona-blog-service.service /etc/systemd/system/"
```

### 5. Configure Systemd Service

```bash
# Reload systemd and enable service
ssh ubuntu@65.1.94.25 << 'EOF'
sudo systemctl daemon-reload
sudo systemctl enable mejona-blog-service
sudo systemctl start mejona-blog-service

# Check status
sudo systemctl status mejona-blog-service
EOF
```

### 6. Configure Nginx

```bash
# Copy nginx configuration
scp nginx/blog-service.conf ubuntu@65.1.94.25:/tmp/
ssh ubuntu@65.1.94.25 << 'EOF'
sudo mv /tmp/blog-service.conf /etc/nginx/sites-available/
sudo ln -sf /etc/nginx/sites-available/blog-service.conf /etc/nginx/sites-enabled/

# Test and reload nginx
sudo nginx -t
sudo systemctl reload nginx
EOF
```

### 7. Configure Firewall

```bash
ssh ubuntu@65.1.94.25 << 'EOF'
# Allow blog service port
sudo ufw allow 8082/tcp comment 'Blog Service'

# Check firewall status
sudo ufw status
EOF
```

## 🐳 Docker Deployment

### Development Environment

```bash
# Start with Docker Compose
docker-compose up -d

# Check services
docker-compose ps

# View logs
docker-compose logs blog-service

# Stop services
docker-compose down
```

### Production Docker

```bash
# Build production image
docker build -t mejona/blog-service:latest .

# Run container
docker run -d \
  --name mejona-blog-service \
  --restart unless-stopped \
  -p 8082:8082 \
  -e GIN_MODE=release \
  -e PORT=8082 \
  --env-file .env \
  mejona/blog-service:latest

# Check container health
docker exec mejona-blog-service curl -f http://localhost:8082/health
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `GIN_MODE` | Gin framework mode | `release` | Yes |
| `PORT` | Service port | `8082` | Yes |
| `DB_HOST` | Database host | `localhost` | Yes |
| `DB_USER` | Database user | `blog_user` | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | `mejona_unified` | Yes |
| `JWT_SECRET` | JWT signing key | - | Yes |
| `ENABLE_SWAGGER` | Enable API docs | `false` | No |
| `LOG_LEVEL` | Log level | `info` | No |

### Database Setup

The service uses the existing `mejona_unified` database. Ensure proper user permissions:

```sql
-- Create database user
CREATE USER 'blog_user'@'localhost' IDENTIFIED BY 'secure_password';

-- Grant permissions
GRANT SELECT, INSERT, UPDATE, DELETE ON mejona_unified.* TO 'blog_user'@'localhost';
GRANT CREATE, ALTER, INDEX ON mejona_unified.* TO 'blog_user'@'localhost';

-- Reload privileges
FLUSH PRIVILEGES;
```

### Redis Setup (Optional)

```bash
# Install Redis (if not present)
sudo apt-get install redis-server

# Configure Redis
sudo sed -i 's/# requirepass foobared/requirepass your_redis_password/' /etc/redis/redis.conf
sudo systemctl restart redis
```

## 🔍 Monitoring & Health Checks

### Health Check Endpoints

```bash
# Basic health check
curl http://65.1.94.25:8082/health

# Deep health check (includes database)
curl http://65.1.94.25:8082/health/deep

# Service status
curl http://65.1.94.25:8082/status

# Kubernetes-style probes
curl http://65.1.94.25:8082/ready   # Readiness probe
curl http://65.1.94.25:8082/alive   # Liveness probe

# Metrics (Prometheus format)
curl http://65.1.94.25:8082/metrics
```

### Nginx Proxy Endpoints

```bash
# Through nginx proxy
curl http://65.1.94.25/blog/health
curl http://65.1.94.25/api/v1/blog/test

# API documentation
curl http://65.1.94.25/blog/docs
```

### Monitoring Scripts

```bash
# Run health check script
./scripts/health-check.sh

# Run comprehensive monitoring
./scripts/monitoring.sh daemon 60

# Generate health report
./scripts/health-check.sh report
```

## 🧪 Testing

### Run Test Suite

```bash
# Run all tests
./scripts/test-runner.sh

# Run specific test types
./scripts/test-runner.sh unit
./scripts/test-runner.sh integration
./scripts/test-runner.sh load

# Run with coverage
./scripts/test-runner.sh coverage
```

### Manual API Testing

```bash
# Test API endpoints
curl -X GET http://65.1.94.25:8082/api/v1/test

# Test with authentication (if implemented)
curl -H "Authorization: Bearer $JWT_TOKEN" \
     http://65.1.94.25:8082/api/v1/protected-endpoint
```

## 📊 CI/CD Pipeline

The service includes automated GitHub Actions workflows:

### Deployment Triggers

- **Push to main/master**: Triggers full deployment pipeline
- **Pull Request**: Runs tests and quality checks
- **Manual Trigger**: Allows manual deployment with options

### Pipeline Stages

1. **Code Quality**: Linting, security scan, formatting
2. **Build & Test**: Unit tests, integration tests, coverage
3. **Docker Build**: Container build and security scan
4. **Staging Deploy**: Deploy to staging environment
5. **Production Deploy**: Deploy to production with verification
6. **Post-Deploy**: Health checks, performance tests, notifications

### Required Secrets

Configure these in GitHub repository secrets:

```
AWS_SSH_PRIVATE_KEY=<SSH private key for EC2>
AWS_USER=ubuntu
PRODUCTION_ENV=<Production environment variables>
SLACK_WEBHOOK=<Slack webhook for notifications>
```

## 🔧 Maintenance

### Log Management

```bash
# View service logs
sudo journalctl -u mejona-blog-service -f

# View nginx logs
sudo tail -f /var/log/nginx/blog-service-access.log
sudo tail -f /var/log/nginx/blog-service-error.log

# View application logs
sudo tail -f /var/log/mejona/blog-service.log
```

### Service Management

```bash
# Service status
sudo systemctl status mejona-blog-service

# Start/stop/restart service
sudo systemctl start mejona-blog-service
sudo systemctl stop mejona-blog-service
sudo systemctl restart mejona-blog-service

# Reload configuration
sudo systemctl reload mejona-blog-service

# View service configuration
sudo systemctl show mejona-blog-service
```

### Database Maintenance

```bash
# Check database connection
mysql -h localhost -u blog_user -p mejona_unified -e "SELECT 1;"

# View service database usage
mysql -h localhost -u blog_user -p mejona_unified -e "SHOW TABLES;"
```

### Performance Monitoring

```bash
# Check resource usage
sudo systemctl status mejona-blog-service
htop -p $(pgrep blog-service)

# Monitor network connections
sudo ss -tlnp | grep :8082
sudo netstat -tulnp | grep :8082

# Check disk usage
df -h /opt/mejona/blog-service
du -sh /var/log/mejona/
```

## 🚨 Troubleshooting

### Common Issues

#### 1. Service Won't Start

**Symptoms:**
- `systemctl status` shows failed state
- Port 8082 not listening

**Solutions:**
```bash
# Check logs for errors
sudo journalctl -u mejona-blog-service -n 50

# Verify configuration
sudo -u bloguser /opt/mejona/blog-service/bin/blog-service --config-check

# Check permissions
ls -la /opt/mejona/blog-service/bin/blog-service
sudo chown bloguser:bloguser /opt/mejona/blog-service/bin/blog-service
sudo chmod +x /opt/mejona/blog-service/bin/blog-service

# Verify environment file
sudo cat /etc/mejona/blog-service.env
```

#### 2. Database Connection Failed

**Symptoms:**
- Deep health check fails
- "Database connection failed" errors

**Solutions:**
```bash
# Test database connection manually
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p $DB_NAME

# Check database user permissions
mysql -u root -p -e "SHOW GRANTS FOR 'blog_user'@'localhost';"

# Verify database exists
mysql -u root -p -e "SHOW DATABASES LIKE 'mejona_unified';"

# Check MySQL service
sudo systemctl status mysql
```

#### 3. Port Already in Use

**Symptoms:**
- "bind: address already in use" error
- Service fails to start

**Solutions:**
```bash
# Find process using port 8082
sudo lsof -i :8082
sudo ss -tlnp | grep :8082

# Kill conflicting process
sudo kill -9 $(sudo lsof -t -i:8082)

# Change service port (if needed)
# Edit /etc/mejona/blog-service.env and set PORT=8083
```

#### 4. High Memory Usage

**Symptoms:**
- Service consuming excessive memory
- System slow or OOM killer activated

**Solutions:**
```bash
# Check memory usage
ps aux | grep blog-service
sudo systemctl show mejona-blog-service --property=MemoryCurrent

# Set memory limits in systemd
sudo systemctl edit mejona-blog-service
# Add:
# [Service]
# MemoryLimit=512M

# Monitor garbage collection
# Add to environment: GODEBUG=gctrace=1
```

#### 5. Nginx Proxy Issues

**Symptoms:**
- `/blog/*` routes return 404
- Upstream connection failed

**Solutions:**
```bash
# Test nginx configuration
sudo nginx -t

# Check upstream health
curl http://localhost:8082/health

# Check nginx error logs
sudo tail -f /var/log/nginx/blog-service-error.log

# Verify nginx site is enabled
ls -la /etc/nginx/sites-enabled/ | grep blog-service

# Restart nginx
sudo systemctl restart nginx
```

#### 6. SSL/TLS Certificate Issues

**Symptoms:**
- HTTPS endpoints fail
- Certificate warnings

**Solutions:**
```bash
# Check certificate validity
openssl x509 -in /etc/nginx/ssl/mejona.in.crt -text -noout

# Verify certificate chain
openssl verify -CAfile /etc/ssl/certs/ca-certificates.crt /etc/nginx/ssl/mejona.in.crt

# Test SSL configuration
sudo nginx -t
```

### Performance Issues

#### Slow Response Times

**Diagnostics:**
```bash
# Monitor response times
curl -w "@curl-format.txt" -o /dev/null -s http://65.1.94.25:8082/health

# Check system load
uptime
top -p $(pgrep blog-service)

# Monitor database queries
# Enable MySQL slow query log
```

**Solutions:**
```bash
# Increase worker connections
# Edit nginx.conf: worker_connections 2048;

# Enable connection pooling
# Verify database connection pool settings

# Add caching headers
# Update nginx configuration with appropriate cache settings
```

### Debugging Tools

#### Service Diagnostics

```bash
# Create service diagnostic script
cat > service-debug.sh << 'EOF'
#!/bin/bash
echo "=== Service Status ==="
sudo systemctl status mejona-blog-service

echo "=== Process Info ==="
ps aux | grep blog-service

echo "=== Network Status ==="
sudo ss -tlnp | grep :8082

echo "=== Recent Logs ==="
sudo journalctl -u mejona-blog-service -n 20 --no-pager

echo "=== Health Check ==="
curl -s http://localhost:8082/health | jq .

echo "=== Resource Usage ==="
df -h /opt/mejona/blog-service
du -sh /var/log/mejona/
EOF

chmod +x service-debug.sh
./service-debug.sh
```

#### Load Testing

```bash
# Install Apache Bench
sudo apt-get install apache2-utils

# Test service under load
ab -n 1000 -c 10 http://65.1.94.25:8082/health

# Monitor during load test
watch -n 1 'sudo systemctl show mejona-blog-service --property=MemoryCurrent,CPUUsage'
```

### Log Analysis

#### Common Log Patterns

```bash
# Find errors in logs
sudo grep -E "(ERROR|FATAL|PANIC)" /var/log/mejona/blog-service.log

# Monitor request patterns
sudo grep "GET\|POST\|PUT\|DELETE" /var/log/nginx/blog-service-access.log | tail -100

# Check response codes
sudo awk '{print $9}' /var/log/nginx/blog-service-access.log | sort | uniq -c
```

#### Log Rotation

```bash
# Configure logrotate for service logs
sudo tee /etc/logrotate.d/mejona-blog-service << 'EOF'
/var/log/mejona/blog-service*.log {
    daily
    missingok
    rotate 30
    compress
    notifempty
    create 0644 bloguser bloguser
    postrotate
        systemctl reload mejona-blog-service > /dev/null 2>&1 || true
    endscript
}
EOF
```

## 📱 API Documentation

### Swagger/OpenAPI

Access API documentation at:
- **Development**: http://localhost:8082/swagger/index.html
- **Production**: http://65.1.94.25:8082/swagger/index.html
- **Nginx Proxy**: http://65.1.94.25/blog/docs

### Core Endpoints

```
GET  /health              - Basic health check
GET  /health/deep         - Comprehensive health check
GET  /status              - Service status
GET  /ready               - Readiness probe
GET  /alive               - Liveness probe
GET  /metrics             - Prometheus metrics
GET  /api/v1/test         - API test endpoint
```

## 🔐 Security

### Best Practices

1. **Environment Variables**: Never commit secrets to repository
2. **TLS/SSL**: Always use HTTPS in production
3. **Firewall**: Restrict access to necessary ports only
4. **User Permissions**: Run service as non-root user
5. **Updates**: Keep dependencies and OS packages updated

### Security Checklist

- [ ] Service runs as non-root user (`bloguser`)
- [ ] Environment file has restricted permissions (600)
- [ ] Database user has minimal required permissions
- [ ] Firewall allows only required ports
- [ ] SSL/TLS certificates are valid and up-to-date
- [ ] Security headers are configured in nginx
- [ ] Regular security scans performed

## 📞 Support

### Getting Help

- **Issues**: Create GitHub issue with deployment logs
- **Documentation**: Check `/docs` directory for additional guides
- **Health Checks**: Use built-in endpoints for status verification

### Contact Information

- **Email**: support@mejona.com
- **GitHub**: https://github.com/MejonaTechnology
- **Website**: https://mejona.in

---

## 📄 Quick Reference

### Key Files & Directories

```
/opt/mejona/blog-service/
├── bin/blog-service           # Main binary
├── logs/                      # Application logs
├── migrations/                # Database migrations
└── docs/                      # Documentation

/etc/mejona/
└── blog-service.env          # Environment configuration

/var/log/mejona/
├── blog-service.log          # Application logs
└── blog-service-error.log    # Error logs

/etc/systemd/system/
└── mejona-blog-service.service  # Systemd service file
```

### Useful Commands

```bash
# Service management
sudo systemctl {start|stop|restart|status} mejona-blog-service

# View logs
sudo journalctl -u mejona-blog-service -f

# Health check
curl http://65.1.94.25:8082/health

# Process information
ps aux | grep blog-service
sudo ss -tlnp | grep :8082

# Deployment
./scripts/deploy-blog-service.sh [--deploy-only|--restart-only]
```

---

**Last Updated**: 2025-01-08  
**Version**: 1.0.0  
**Service Port**: 8082  
**Deployment Target**: AWS EC2 (65.1.94.25)
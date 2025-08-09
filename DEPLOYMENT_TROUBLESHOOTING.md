# Blog CRM Service - Deployment Troubleshooting Guide

## Overview

This guide provides troubleshooting steps for common deployment issues with the Blog CRM Management Microservice.

## Quick Status Check

```bash
# Check service status
sudo systemctl status blog-crm-service

# Check if service is listening on port 8082
sudo ss -tlnp | grep 8082

# Test health endpoint
curl http://localhost:8082/health

# Test external access
curl http://65.1.94.25:8082/health
```

## Common Issues and Solutions

### 1. Service Failed to Start

**Symptoms:**
- `systemctl status blog-crm-service` shows "failed" or "inactive"
- Service logs show startup errors

**Debugging Steps:**
```bash
# Check service logs
sudo journalctl -u blog-crm-service -n 50

# Check if binary exists and is executable
ls -la /opt/mejona/blog-crm-service/blog-crm-service

# Test manual startup
cd /opt/mejona/blog-crm-service
sudo -u ubuntu ./blog-crm-service
```

**Common Causes & Solutions:**

#### Missing Environment Variables
```bash
# Check environment file
cat /opt/mejona/blog-crm-service/.env

# Required variables:
# DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET
```

#### Database Connection Issues
```bash
# Test database connection
mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME -e "SELECT 1;"
```

#### Port Already in Use
```bash
# Check what's using port 8082
sudo lsof -i :8082

# Kill conflicting process if necessary
sudo kill -9 <PID>
```

#### Permission Issues
```bash
# Fix ownership and permissions
sudo chown -R ubuntu:ubuntu /opt/mejona/blog-crm-service
chmod +x /opt/mejona/blog-crm-service/blog-crm-service
```

### 2. Build Failures

**Symptoms:**
- GitHub Actions build step fails
- Local build produces errors

**Debugging Steps:**
```bash
# Check Go version
go version

# Clean and rebuild
go clean -cache
go mod download
go mod verify
go build -v cmd/server/main.go
```

**Common Causes & Solutions:**

#### Go Version Mismatch
- Ensure workflow uses Go 1.23 (matches go.mod)
- Update local Go installation if necessary

#### Missing Dependencies
```bash
go mod tidy
go mod download
```

#### Build Environment Issues
```bash
# Set build environment
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
```

### 3. Health Check Failures

**Symptoms:**
- Health endpoint returns 500 or timeouts
- Service appears running but health checks fail

**Debugging Steps:**
```bash
# Check service logs for errors
sudo journalctl -u blog-crm-service -f

# Test health endpoint with verbose output
curl -v http://localhost:8082/health

# Check database connectivity
curl http://localhost:8082/health/deep
```

**Solutions:**
- Verify database connection parameters
- Check network connectivity to database
- Ensure all required environment variables are set

### 4. CI/CD Pipeline Failures

**Symptoms:**
- GitHub Actions workflow fails
- Deployment steps timeout or error

**Common Issues:**

#### SSH Connection Issues
- Verify `PRODUCTION_HOST` secret is correct (65.1.94.25)
- Check `EC2_SSH_KEY` is valid and has proper permissions
- Ensure `EC2_USERNAME` is set to "ubuntu"

#### Missing Secrets
Required GitHub secrets:
- `PRODUCTION_HOST`: 65.1.94.25
- `EC2_USERNAME`: ubuntu
- `EC2_SSH_KEY`: Private SSH key for EC2 access
- `DB_HOST`: Database host
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (mejona_unified)
- `JWT_SECRET`: JWT signing secret

#### Timeout Issues
- Increase timeout values in workflow
- Check EC2 instance resources and availability
- Verify network connectivity

### 5. External Access Issues

**Symptoms:**
- Service works on localhost but not externally
- Curl from outside server fails

**Debugging Steps:**
```bash
# Check if service binds to all interfaces
sudo netstat -tlnp | grep 8082

# Check iptables/firewall rules
sudo iptables -L | grep 8082

# Test from another machine
curl http://65.1.94.25:8082/health
```

**Solutions:**
- Ensure service binds to 0.0.0.0:8082, not just localhost
- Configure security groups in AWS to allow port 8082
- Update firewall rules if necessary

## Monitoring and Maintenance

### Log Files
```bash
# Service logs (systemd journal)
sudo journalctl -u blog-crm-service -f

# Application logs (if configured)
tail -f /var/log/mejona/blog-crm-service.log

# Health check logs
tail -f /var/log/mejona/blog-crm-service-health.log
```

### Health Monitoring
```bash
# Run comprehensive health check
./scripts/health-check.sh comprehensive

# Continuous monitoring
./scripts/health-check.sh monitor 60

# Generate health report
./scripts/health-check.sh report
```

### Manual Recovery
```bash
# Stop service
sudo systemctl stop blog-crm-service

# Update code
cd /opt/mejona/blog-crm-service
sudo git pull origin main

# Rebuild
sudo -u ubuntu /usr/local/go/bin/go build -o blog-crm-service cmd/server/main.go

# Start service
sudo systemctl start blog-crm-service

# Verify
curl http://localhost:8082/health
```

## Performance Optimization

### Resource Monitoring
```bash
# Check memory usage
ps aux | grep blog-crm-service

# Check CPU usage
top -p $(pgrep blog-crm-service)

# Check disk usage
df -h /opt/mejona/blog-crm-service
```

### Database Optimization
```bash
# Check database connections
mysql -e "SHOW PROCESSLIST;" -u $DB_USER -p$DB_PASSWORD $DB_NAME

# Monitor slow queries
mysql -e "SHOW VARIABLES LIKE 'slow_query_log';" -u $DB_USER -p$DB_PASSWORD $DB_NAME
```

## Emergency Procedures

### Service Recovery
```bash
# Automated recovery attempt
./scripts/health-check.sh recover

# Manual recovery
sudo systemctl restart blog-crm-service
sleep 10
curl http://localhost:8082/health
```

### Rollback Deployment
```bash
# Restore from backup
cd /opt/mejona/blog-crm-service
sudo cp blog-crm-service.backup.* blog-crm-service
sudo systemctl restart blog-crm-service
```

### Contact Support
If issues persist:
1. Generate health report: `./scripts/health-check.sh report`
2. Collect service logs: `sudo journalctl -u blog-crm-service -n 100`
3. Document reproduction steps
4. Contact system administrator with collected information

## Preventive Measures

### Regular Monitoring
- Set up cron job for health checks
- Monitor disk space and memory usage
- Review logs regularly for warnings

### Backup Strategy
- Automated database backups
- Service configuration backups
- Binary backups before deployments

### Testing
- Test deployments in staging environment
- Validate health endpoints after each deployment
- Monitor resource usage trends

## Useful Commands Reference

```bash
# Service Management
sudo systemctl status blog-crm-service
sudo systemctl start blog-crm-service
sudo systemctl stop blog-crm-service
sudo systemctl restart blog-crm-service
sudo systemctl reload blog-crm-service

# Logs
sudo journalctl -u blog-crm-service -f
sudo journalctl -u blog-crm-service -n 50
sudo journalctl -u blog-crm-service --since "1 hour ago"

# Health Checks
curl http://localhost:8082/health
curl http://localhost:8082/health/deep
curl http://localhost:8082/api/v1/test
curl http://localhost:8082/metrics

# Network
sudo ss -tlnp | grep 8082
sudo netstat -tlnp | grep 8082
sudo lsof -i :8082

# Process Management
ps aux | grep blog-crm-service
pgrep blog-crm-service
pkill blog-crm-service
```

---

**Last Updated:** $(date '+%Y-%m-%d %H:%M:%S')  
**Version:** 1.0  
**Service:** Blog CRM Management Microservice  
**Port:** 8082
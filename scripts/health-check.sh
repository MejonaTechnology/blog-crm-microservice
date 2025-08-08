#!/bin/bash
# =============================================================================
# Blog CRM Microservice - Health Check & Monitoring Script
# Mejona Technology LLP
# =============================================================================

set -euo pipefail

# =============================================================================
# CONFIGURATION
# =============================================================================
readonly SERVICE_NAME="blog-service"
readonly SERVICE_PORT="8082"
readonly SYSTEMD_SERVICE="mejona-blog-service"
readonly LOG_FILE="/var/log/mejona/blog-service.log"
readonly ERROR_LOG="/var/log/mejona/blog-service-error.log"
readonly HEALTH_LOG="/var/log/mejona/blog-service-health.log"
readonly PID_FILE="/var/run/mejona/blog-service.pid"

# Colors
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'

# Health check endpoints
readonly BASE_URL="http://localhost:$SERVICE_PORT"
readonly HEALTH_ENDPOINT="$BASE_URL/health"
readonly DEEP_HEALTH_ENDPOINT="$BASE_URL/health/deep"
readonly METRICS_ENDPOINT="$BASE_URL/metrics"
readonly API_TEST_ENDPOINT="$BASE_URL/api/v1/test"

# Alert thresholds
readonly MAX_RESPONSE_TIME=5000    # milliseconds
readonly MAX_CPU_PERCENT=80
readonly MAX_MEMORY_PERCENT=75
readonly MAX_DISK_PERCENT=85

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

log_health() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> "$HEALTH_LOG"
}

# =============================================================================
# HEALTH CHECK FUNCTIONS
# =============================================================================

check_service_status() {
    local status
    if systemctl is-active "$SYSTEMD_SERVICE" >/dev/null 2>&1; then
        status=$(systemctl is-active "$SYSTEMD_SERVICE")
        if [[ "$status" == "active" ]]; then
            success "Service is active and running"
            log_health "SERVICE_STATUS: ACTIVE"
            return 0
        else
            error "Service status: $status"
            log_health "SERVICE_STATUS: $status"
            return 1
        fi
    else
        error "Service is not running"
        log_health "SERVICE_STATUS: INACTIVE"
        return 1
    fi
}

check_port_listening() {
    if ss -tlnp | grep -q ":$SERVICE_PORT "; then
        success "Service is listening on port $SERVICE_PORT"
        log_health "PORT_CHECK: LISTENING"
        return 0
    else
        error "Service is not listening on port $SERVICE_PORT"
        log_health "PORT_CHECK: NOT_LISTENING"
        return 1
    fi
}

check_health_endpoint() {
    local response_code
    local response_time
    local start_time
    local end_time
    
    start_time=$(date +%s%3N)
    
    if response_code=$(curl -s -w "%{http_code}" -o /tmp/health_response.json "$HEALTH_ENDPOINT" --connect-timeout 10 --max-time 30); then
        end_time=$(date +%s%3N)
        response_time=$((end_time - start_time))
        
        if [[ "$response_code" == "200" ]]; then
            success "Health endpoint responded with 200 OK (${response_time}ms)"
            log_health "HEALTH_ENDPOINT: OK - ${response_time}ms"
            
            if [[ $response_time -gt $MAX_RESPONSE_TIME ]]; then
                warn "Response time ${response_time}ms exceeds threshold ${MAX_RESPONSE_TIME}ms"
                log_health "HEALTH_ENDPOINT: SLOW_RESPONSE - ${response_time}ms"
            fi
            
            return 0
        else
            error "Health endpoint responded with $response_code"
            log_health "HEALTH_ENDPOINT: ERROR - HTTP $response_code"
            return 1
        fi
    else
        error "Health endpoint is unreachable"
        log_health "HEALTH_ENDPOINT: UNREACHABLE"
        return 1
    fi
}

check_deep_health() {
    local response_code
    
    if response_code=$(curl -s -w "%{http_code}" -o /tmp/deep_health_response.json "$DEEP_HEALTH_ENDPOINT" --connect-timeout 10 --max-time 30); then
        if [[ "$response_code" == "200" ]]; then
            success "Deep health check passed"
            log_health "DEEP_HEALTH: PASSED"
            
            # Parse and display database and dependencies status
            if command -v jq >/dev/null 2>&1 && [[ -f /tmp/deep_health_response.json ]]; then
                local db_status
                db_status=$(jq -r '.database.status // "unknown"' /tmp/deep_health_response.json 2>/dev/null || echo "unknown")
                info "Database status: $db_status"
                log_health "DATABASE_STATUS: $db_status"
            fi
            
            return 0
        else
            error "Deep health check failed with HTTP $response_code"
            log_health "DEEP_HEALTH: FAILED - HTTP $response_code"
            return 1
        fi
    else
        error "Deep health check is unreachable"
        log_health "DEEP_HEALTH: UNREACHABLE"
        return 1
    fi
}

check_api_endpoint() {
    local response_code
    
    if response_code=$(curl -s -w "%{http_code}" -o /tmp/api_test_response.json "$API_TEST_ENDPOINT" --connect-timeout 10 --max-time 30); then
        if [[ "$response_code" == "200" ]]; then
            success "API test endpoint is working"
            log_health "API_TEST: OK"
            return 0
        else
            error "API test endpoint failed with HTTP $response_code"
            log_health "API_TEST: FAILED - HTTP $response_code"
            return 1
        fi
    else
        error "API test endpoint is unreachable"
        log_health "API_TEST: UNREACHABLE"
        return 1
    fi
}

# =============================================================================
# SYSTEM RESOURCE MONITORING
# =============================================================================

check_system_resources() {
    local pid
    local cpu_percent=0
    local memory_percent=0
    local disk_percent
    
    # Get service PID
    if pid=$(systemctl show "$SYSTEMD_SERVICE" --property=MainPID --value) && [[ "$pid" != "0" ]]; then
        # Get CPU and memory usage
        if command -v ps >/dev/null 2>&1; then
            local ps_output
            if ps_output=$(ps -p "$pid" -o pid,ppid,%cpu,%mem,rss,vsz,cmd --no-headers 2>/dev/null); then
                cpu_percent=$(echo "$ps_output" | awk '{print $3}')
                memory_percent=$(echo "$ps_output" | awk '{print $4}')
                local rss=$(echo "$ps_output" | awk '{print $5}')
                local vsz=$(echo "$ps_output" | awk '{print $6}')
                
                info "Process PID: $pid"
                info "CPU Usage: ${cpu_percent}%"
                info "Memory Usage: ${memory_percent}% (RSS: ${rss}KB, VSZ: ${vsz}KB)"
                
                log_health "RESOURCE_CHECK: PID=$pid CPU=${cpu_percent}% MEM=${memory_percent}%"
                
                # Check thresholds
                if (( $(echo "$cpu_percent > $MAX_CPU_PERCENT" | bc -l) )); then
                    warn "CPU usage ${cpu_percent}% exceeds threshold ${MAX_CPU_PERCENT}%"
                    log_health "ALERT: HIGH_CPU - ${cpu_percent}%"
                fi
                
                if (( $(echo "$memory_percent > $MAX_MEMORY_PERCENT" | bc -l) )); then
                    warn "Memory usage ${memory_percent}% exceeds threshold ${MAX_MEMORY_PERCENT}%"
                    log_health "ALERT: HIGH_MEMORY - ${memory_percent}%"
                fi
            fi
        fi
    else
        warn "Could not determine service PID"
        log_health "RESOURCE_CHECK: PID_NOT_FOUND"
    fi
    
    # Check disk usage
    if command -v df >/dev/null 2>&1; then
        disk_percent=$(df /opt/mejona/blog-service | tail -1 | awk '{print $5}' | sed 's/%//')
        info "Disk Usage: ${disk_percent}%"
        log_health "DISK_USAGE: ${disk_percent}%"
        
        if [[ $disk_percent -gt $MAX_DISK_PERCENT ]]; then
            warn "Disk usage ${disk_percent}% exceeds threshold ${MAX_DISK_PERCENT}%"
            log_health "ALERT: HIGH_DISK_USAGE - ${disk_percent}%"
        fi
    fi
}

check_log_files() {
    local log_files=("$LOG_FILE" "$ERROR_LOG")
    local recent_errors=0
    
    for log_file in "${log_files[@]}"; do
        if [[ -f "$log_file" ]]; then
            local size
            size=$(ls -lh "$log_file" | awk '{print $5}')
            info "Log file $log_file: $size"
            
            # Check for recent errors (last 5 minutes)
            if [[ "$log_file" == "$ERROR_LOG" ]]; then
                recent_errors=$(find "$log_file" -newermt "5 minutes ago" -exec grep -c "ERROR\|FATAL\|PANIC" {} \; 2>/dev/null || echo 0)
                if [[ $recent_errors -gt 0 ]]; then
                    warn "Found $recent_errors recent error(s) in $log_file"
                    log_health "LOG_CHECK: RECENT_ERRORS - $recent_errors"
                fi
            fi
        else
            warn "Log file $log_file does not exist"
            log_health "LOG_CHECK: MISSING - $log_file"
        fi
    done
}

check_database_connection() {
    local db_host="${DB_HOST:-localhost}"
    local db_port="${DB_PORT:-3306}"
    local db_user="${DB_USER:-blog_user}"
    local db_name="${DB_NAME:-mejona_unified}"
    
    if command -v mysql >/dev/null 2>&1; then
        if mysql -h "$db_host" -P "$db_port" -u "$db_user" -e "SELECT 1;" "$db_name" >/dev/null 2>&1; then
            success "Database connection successful"
            log_health "DATABASE: CONNECTED"
            return 0
        else
            error "Database connection failed"
            log_health "DATABASE: CONNECTION_FAILED"
            return 1
        fi
    else
        warn "MySQL client not available for database check"
        log_health "DATABASE: CLIENT_NOT_AVAILABLE"
        return 1
    fi
}

# =============================================================================
# AUTOMATED RECOVERY
# =============================================================================

attempt_service_recovery() {
    warn "Attempting service recovery..."
    log_health "RECOVERY: ATTEMPTING_RESTART"
    
    # Try to restart the service
    if systemctl restart "$SYSTEMD_SERVICE"; then
        sleep 10  # Wait for service to start
        
        # Check if recovery was successful
        if check_service_status && check_health_endpoint; then
            success "Service recovery successful"
            log_health "RECOVERY: SUCCESS"
            return 0
        else
            error "Service recovery failed"
            log_health "RECOVERY: FAILED"
            return 1
        fi
    else
        error "Failed to restart service"
        log_health "RECOVERY: RESTART_FAILED"
        return 1
    fi
}

# =============================================================================
# REPORTING FUNCTIONS
# =============================================================================

generate_health_report() {
    local report_file="/tmp/blog-service-health-report-$(date +%Y%m%d-%H%M%S).txt"
    
    {
        echo "=============================================="
        echo "Mejona Blog Service Health Report"
        echo "Generated: $(date '+%Y-%m-%d %H:%M:%S')"
        echo "=============================================="
        echo
        
        echo "SERVICE STATUS:"
        systemctl status "$SYSTEMD_SERVICE" --no-pager || true
        echo
        
        echo "NETWORK STATUS:"
        ss -tlnp | grep ":$SERVICE_PORT" || echo "Port $SERVICE_PORT not listening"
        echo
        
        echo "RESOURCE USAGE:"
        if pid=$(systemctl show "$SYSTEMD_SERVICE" --property=MainPID --value) && [[ "$pid" != "0" ]]; then
            ps -p "$pid" -o pid,ppid,%cpu,%mem,rss,vsz,cmd || true
        fi
        echo
        
        echo "DISK USAGE:"
        df -h /opt/mejona/blog-service || true
        echo
        
        echo "RECENT LOG ENTRIES:"
        tail -20 "$LOG_FILE" 2>/dev/null || echo "Log file not accessible"
        echo
        
        echo "RECENT ERROR LOG ENTRIES:"
        tail -10 "$ERROR_LOG" 2>/dev/null || echo "Error log file not accessible"
        echo
        
        echo "HEALTH CHECK RESULTS:"
        curl -s "$HEALTH_ENDPOINT" | jq . 2>/dev/null || curl -s "$HEALTH_ENDPOINT" || echo "Health endpoint not accessible"
        
    } > "$report_file"
    
    info "Health report generated: $report_file"
    echo "$report_file"
}

# =============================================================================
# MAIN FUNCTIONS
# =============================================================================

run_basic_health_check() {
    local checks_passed=0
    local total_checks=4
    
    log "Running basic health check..."
    
    if check_service_status; then
        ((checks_passed++))
    fi
    
    if check_port_listening; then
        ((checks_passed++))
    fi
    
    if check_health_endpoint; then
        ((checks_passed++))
    fi
    
    if check_api_endpoint; then
        ((checks_passed++))
    fi
    
    info "Basic health check: $checks_passed/$total_checks checks passed"
    log_health "BASIC_CHECK: $checks_passed/$total_checks PASSED"
    
    return $((total_checks - checks_passed))
}

run_comprehensive_health_check() {
    local checks_passed=0
    local total_checks=7
    
    log "Running comprehensive health check..."
    
    # Load environment variables if available
    if [[ -f /etc/mejona/blog-service.env ]]; then
        set -a
        source /etc/mejona/blog-service.env
        set +a
    fi
    
    if check_service_status; then
        ((checks_passed++))
    fi
    
    if check_port_listening; then
        ((checks_passed++))
    fi
    
    if check_health_endpoint; then
        ((checks_passed++))
    fi
    
    if check_deep_health; then
        ((checks_passed++))
    fi
    
    if check_api_endpoint; then
        ((checks_passed++))
    fi
    
    check_system_resources
    ((checks_passed++))  # Always count as passed for monitoring purposes
    
    check_log_files
    
    if check_database_connection; then
        ((checks_passed++))
    fi
    
    info "Comprehensive health check: $checks_passed/$total_checks checks passed"
    log_health "COMPREHENSIVE_CHECK: $checks_passed/$total_checks PASSED"
    
    return $((total_checks - checks_passed))
}

monitor_service() {
    local interval="${1:-60}"
    local max_failures="${2:-3}"
    local failure_count=0
    
    log "Starting continuous monitoring (interval: ${interval}s, max failures: $max_failures)"
    
    while true; do
        if run_basic_health_check; then
            failure_count=0
            log_health "MONITOR: HEALTHY"
        else
            ((failure_count++))
            warn "Health check failed (failure $failure_count/$max_failures)"
            log_health "MONITOR: UNHEALTHY - failure $failure_count/$max_failures"
            
            if [[ $failure_count -ge $max_failures ]]; then
                error "Maximum failures reached, attempting recovery..."
                if attempt_service_recovery; then
                    failure_count=0
                else
                    error "Recovery failed, generating report and exiting..."
                    generate_health_report
                    exit 1
                fi
            fi
        fi
        
        sleep "$interval"
    done
}

show_help() {
    cat << EOF
Usage: $0 [COMMAND] [OPTIONS]

Health check and monitoring script for Mejona Blog Service

COMMANDS:
    check                   Run basic health check (default)
    comprehensive           Run comprehensive health check
    monitor [interval]      Start continuous monitoring
    report                  Generate health report
    recover                 Attempt service recovery
    status                  Show service status
    logs [lines]           Show recent log entries

OPTIONS:
    -h, --help             Show this help message
    -v, --verbose          Verbose output
    -q, --quiet            Quiet output (errors only)

EXAMPLES:
    $0                     Run basic health check
    $0 comprehensive       Run comprehensive health check
    $0 monitor 30          Monitor every 30 seconds
    $0 report              Generate health report
    $0 logs 50             Show last 50 log lines

EOF
}

main() {
    local command="${1:-check}"
    local verbose=false
    local quiet=false
    
    # Ensure log directory exists
    mkdir -p "$(dirname "$HEALTH_LOG")"
    
    case "$command" in
        -h|--help|help)
            show_help
            exit 0
            ;;
        check|basic)
            run_basic_health_check
            exit $?
            ;;
        comprehensive|full)
            run_comprehensive_health_check
            exit $?
            ;;
        monitor)
            local interval="${2:-60}"
            monitor_service "$interval"
            ;;
        report)
            report_file=$(generate_health_report)
            cat "$report_file"
            ;;
        recover|recovery)
            attempt_service_recovery
            exit $?
            ;;
        status)
            systemctl status "$SYSTEMD_SERVICE" --no-pager
            ;;
        logs)
            local lines="${2:-20}"
            echo "=== Service Logs (last $lines lines) ==="
            tail -n "$lines" "$LOG_FILE" 2>/dev/null || echo "Log file not accessible"
            echo
            echo "=== Error Logs (last $lines lines) ==="
            tail -n "$lines" "$ERROR_LOG" 2>/dev/null || echo "Error log file not accessible"
            ;;
        *)
            error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
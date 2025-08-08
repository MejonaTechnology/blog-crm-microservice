#!/bin/bash
# =============================================================================
# Blog CRM Microservice - Advanced Monitoring & Alerting Script
# Mejona Technology LLP
# =============================================================================

set -euo pipefail

# =============================================================================
# CONFIGURATION
# =============================================================================
readonly SERVICE_NAME="blog-service"
readonly SERVICE_PORT="8082"
readonly SYSTEMD_SERVICE="mejona-blog-service"
readonly METRICS_DIR="/var/lib/mejona/metrics"
readonly ALERT_DIR="/var/lib/mejona/alerts"
readonly MONITORING_LOG="/var/log/mejona/blog-service-monitoring.log"

# Metrics collection endpoints
readonly BASE_URL="http://localhost:$SERVICE_PORT"
readonly METRICS_ENDPOINT="$BASE_URL/metrics"
readonly HEALTH_ENDPOINT="$BASE_URL/health"

# Alert thresholds
readonly CPU_WARNING=60
readonly CPU_CRITICAL=80
readonly MEMORY_WARNING=70
readonly MEMORY_CRITICAL=85
readonly DISK_WARNING=80
readonly DISK_CRITICAL=90
readonly RESPONSE_TIME_WARNING=2000
readonly RESPONSE_TIME_CRITICAL=5000
readonly ERROR_RATE_WARNING=5
readonly ERROR_RATE_CRITICAL=10

# Alert notification settings
readonly ALERT_EMAIL="${ALERT_EMAIL:-admin@mejona.com}"
readonly SLACK_WEBHOOK="${SLACK_WEBHOOK:-}"
readonly ALERT_COOLDOWN=300  # 5 minutes

# Colors
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'

# =============================================================================
# UTILITY FUNCTIONS
# =============================================================================

log() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" >> "$MONITORING_LOG"
}

warn() {
    echo -e "${YELLOW}[$(date '+%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] WARNING: $1" >> "$MONITORING_LOG"
}

error() {
    echo -e "${RED}[$(date '+%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1" >&2
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] ERROR: $1" >> "$MONITORING_LOG"
}

success() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')] SUCCESS:${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] SUCCESS: $1" >> "$MONITORING_LOG"
}

info() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')] INFO:${NC} $1"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] INFO: $1" >> "$MONITORING_LOG"
}

# =============================================================================
# METRICS COLLECTION
# =============================================================================

collect_system_metrics() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    local metrics_file="$METRICS_DIR/system_metrics_$(date +%Y%m%d).json"
    
    # Get service PID
    local pid
    if ! pid=$(systemctl show "$SYSTEMD_SERVICE" --property=MainPID --value) || [[ "$pid" == "0" ]]; then
        warn "Could not determine service PID"
        return 1
    fi
    
    # Collect CPU and memory metrics
    local cpu_percent=0
    local memory_percent=0
    local rss_kb=0
    local vsz_kb=0
    
    if ps_output=$(ps -p "$pid" -o %cpu,%mem,rss,vsz --no-headers 2>/dev/null); then
        cpu_percent=$(echo "$ps_output" | awk '{print $1}')
        memory_percent=$(echo "$ps_output" | awk '{print $2}')
        rss_kb=$(echo "$ps_output" | awk '{print $3}')
        vsz_kb=$(echo "$ps_output" | awk '{print $4}')
    fi
    
    # Collect disk metrics
    local disk_percent=0
    if df_output=$(df /opt/mejona/blog-service 2>/dev/null); then
        disk_percent=$(echo "$df_output" | tail -1 | awk '{print $5}' | sed 's/%//')
    fi
    
    # Collect network metrics
    local connections=0
    if ss_output=$(ss -tn state established "( dport = :$SERVICE_PORT or sport = :$SERVICE_PORT )" 2>/dev/null); then
        connections=$(echo "$ss_output" | wc -l)
        ((connections--)) # Subtract header line
    fi
    
    # Create JSON metrics
    local metrics_json
    metrics_json=$(cat << EOF
{
    "timestamp": "$timestamp",
    "service": "$SERVICE_NAME",
    "pid": $pid,
    "cpu_percent": $cpu_percent,
    "memory_percent": $memory_percent,
    "rss_kb": $rss_kb,
    "vsz_kb": $vsz_kb,
    "disk_percent": $disk_percent,
    "active_connections": $connections
}
EOF
)
    
    # Append to metrics file
    echo "$metrics_json" >> "$metrics_file"
    
    # Return metrics for threshold checking
    echo "$cpu_percent|$memory_percent|$disk_percent|$connections"
}

collect_application_metrics() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    local metrics_file="$METRICS_DIR/app_metrics_$(date +%Y%m%d).json"
    
    # Test health endpoint response time
    local start_time=$(date +%s%3N)
    local response_code
    local response_time=0
    
    if response_code=$(curl -s -w "%{http_code}" -o /dev/null "$HEALTH_ENDPOINT" --connect-timeout 10 --max-time 30); then
        local end_time=$(date +%s%3N)
        response_time=$((end_time - start_time))
    else
        response_code=0
        response_time=30000  # Timeout
    fi
    
    # Check if metrics endpoint is available
    local metrics_available=false
    local request_count=0
    local error_count=0
    local success_rate=100
    
    if curl -sf "$METRICS_ENDPOINT" >/dev/null 2>&1; then
        metrics_available=true
        # Parse metrics if available (Prometheus format)
        if metrics_data=$(curl -s "$METRICS_ENDPOINT" 2>/dev/null); then
            request_count=$(echo "$metrics_data" | grep -E '^http_requests_total' | awk '{sum+=$2} END {print sum+0}')
            error_count=$(echo "$metrics_data" | grep -E '^http_requests_total.*5[0-9][0-9]' | awk '{sum+=$2} END {print sum+0}')
            if [[ $request_count -gt 0 ]]; then
                success_rate=$(( (request_count - error_count) * 100 / request_count ))
            fi
        fi
    fi
    
    # Create JSON metrics
    local app_metrics_json
    app_metrics_json=$(cat << EOF
{
    "timestamp": "$timestamp",
    "service": "$SERVICE_NAME",
    "health_response_code": $response_code,
    "response_time_ms": $response_time,
    "metrics_available": $metrics_available,
    "total_requests": $request_count,
    "error_requests": $error_count,
    "success_rate_percent": $success_rate
}
EOF
)
    
    # Append to metrics file
    echo "$app_metrics_json" >> "$metrics_file"
    
    # Return metrics for threshold checking
    echo "$response_code|$response_time|$((100 - success_rate))"
}

collect_log_metrics() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    local metrics_file="$METRICS_DIR/log_metrics_$(date +%Y%m%d).json"
    local log_file="/var/log/mejona/blog-service.log"
    local error_log="/var/log/mejona/blog-service-error.log"
    
    # Count log entries in the last 5 minutes
    local info_count=0
    local warn_count=0
    local error_count=0
    
    if [[ -f "$log_file" ]]; then
        info_count=$(grep -c "INFO" "$log_file" 2>/dev/null | tail -100 || echo 0)
        warn_count=$(grep -c "WARN\|WARNING" "$log_file" 2>/dev/null | tail -100 || echo 0)
    fi
    
    if [[ -f "$error_log" ]]; then
        error_count=$(grep -c "ERROR\|FATAL\|PANIC" "$error_log" 2>/dev/null | tail -100 || echo 0)
    fi
    
    # Check log file sizes
    local log_size=0
    local error_log_size=0
    
    if [[ -f "$log_file" ]]; then
        log_size=$(stat -c%s "$log_file" 2>/dev/null || echo 0)
    fi
    
    if [[ -f "$error_log" ]]; then
        error_log_size=$(stat -c%s "$error_log" 2>/dev/null || echo 0)
    fi
    
    # Create JSON metrics
    local log_metrics_json
    log_metrics_json=$(cat << EOF
{
    "timestamp": "$timestamp",
    "service": "$SERVICE_NAME",
    "info_count": $info_count,
    "warn_count": $warn_count,
    "error_count": $error_count,
    "log_size_bytes": $log_size,
    "error_log_size_bytes": $error_log_size
}
EOF
)
    
    # Append to metrics file
    echo "$log_metrics_json" >> "$metrics_file"
    
    # Return error count for alerting
    echo "$error_count"
}

# =============================================================================
# ALERTING SYSTEM
# =============================================================================

should_send_alert() {
    local alert_type="$1"
    local cooldown_file="$ALERT_DIR/${alert_type}_last_sent"
    local current_time=$(date +%s)
    
    if [[ -f "$cooldown_file" ]]; then
        local last_sent
        last_sent=$(cat "$cooldown_file")
        local time_diff=$((current_time - last_sent))
        
        if [[ $time_diff -lt $ALERT_COOLDOWN ]]; then
            return 1  # Still in cooldown
        fi
    fi
    
    echo "$current_time" > "$cooldown_file"
    return 0  # Can send alert
}

send_email_alert() {
    local subject="$1"
    local message="$2"
    local severity="${3:-WARNING}"
    
    if [[ -z "$ALERT_EMAIL" ]]; then
        return 0
    fi
    
    # Create email content
    local email_content
    email_content=$(cat << EOF
Subject: [Mejona Blog Service] $severity: $subject
To: $ALERT_EMAIL
Content-Type: text/html

<html>
<body>
<h2>Mejona Blog Service Alert</h2>
<p><strong>Severity:</strong> $severity</p>
<p><strong>Service:</strong> $SERVICE_NAME</p>
<p><strong>Time:</strong> $(date)</p>
<p><strong>Server:</strong> $(hostname)</p>
<hr>
<p>$message</p>
<hr>
<p><em>This is an automated alert from the Mejona monitoring system.</em></p>
</body>
</html>
EOF
)
    
    # Send email using sendmail if available
    if command -v sendmail >/dev/null 2>&1; then
        echo "$email_content" | sendmail "$ALERT_EMAIL"
        log "Email alert sent to $ALERT_EMAIL"
    elif command -v mail >/dev/null 2>&1; then
        echo "$message" | mail -s "[Mejona Blog Service] $severity: $subject" "$ALERT_EMAIL"
        log "Mail alert sent to $ALERT_EMAIL"
    else
        warn "No mail command available for email alerts"
    fi
}

send_slack_alert() {
    local message="$1"
    local severity="${2:-warning}"
    
    if [[ -z "$SLACK_WEBHOOK" ]]; then
        return 0
    fi
    
    local color="warning"
    case "$severity" in
        "CRITICAL") color="danger" ;;
        "WARNING") color="warning" ;;
        "INFO") color="good" ;;
    esac
    
    local slack_payload
    slack_payload=$(cat << EOF
{
    "attachments": [
        {
            "color": "$color",
            "title": "Mejona Blog Service Alert",
            "fields": [
                {
                    "title": "Severity",
                    "value": "$severity",
                    "short": true
                },
                {
                    "title": "Service",
                    "value": "$SERVICE_NAME",
                    "short": true
                },
                {
                    "title": "Server",
                    "value": "$(hostname)",
                    "short": true
                },
                {
                    "title": "Time",
                    "value": "$(date)",
                    "short": true
                },
                {
                    "title": "Details",
                    "value": "$message",
                    "short": false
                }
            ]
        }
    ]
}
EOF
)
    
    if curl -X POST -H 'Content-type: application/json' --data "$slack_payload" "$SLACK_WEBHOOK" >/dev/null 2>&1; then
        log "Slack alert sent"
    else
        warn "Failed to send Slack alert"
    fi
}

trigger_alert() {
    local alert_type="$1"
    local severity="$2"
    local subject="$3"
    local message="$4"
    
    if should_send_alert "$alert_type"; then
        log "Triggering $severity alert: $subject"
        
        # Log the alert
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] $severity ALERT: $subject - $message" >> "$ALERT_DIR/alerts.log"
        
        # Send notifications
        send_email_alert "$subject" "$message" "$severity"
        send_slack_alert "$message" "$severity"
    else
        info "Alert $alert_type is in cooldown period"
    fi
}

# =============================================================================
# THRESHOLD CHECKING
# =============================================================================

check_system_thresholds() {
    local metrics="$1"
    IFS='|' read -r cpu_percent memory_percent disk_percent connections <<< "$metrics"
    
    # CPU alerts
    if (( $(echo "$cpu_percent > $CPU_CRITICAL" | bc -l) )); then
        trigger_alert "high_cpu" "CRITICAL" "High CPU Usage" "CPU usage is at ${cpu_percent}% (critical threshold: ${CPU_CRITICAL}%)"
    elif (( $(echo "$cpu_percent > $CPU_WARNING" | bc -l) )); then
        trigger_alert "high_cpu" "WARNING" "Elevated CPU Usage" "CPU usage is at ${cpu_percent}% (warning threshold: ${CPU_WARNING}%)"
    fi
    
    # Memory alerts
    if (( $(echo "$memory_percent > $MEMORY_CRITICAL" | bc -l) )); then
        trigger_alert "high_memory" "CRITICAL" "High Memory Usage" "Memory usage is at ${memory_percent}% (critical threshold: ${MEMORY_CRITICAL}%)"
    elif (( $(echo "$memory_percent > $MEMORY_WARNING" | bc -l) )); then
        trigger_alert "high_memory" "WARNING" "Elevated Memory Usage" "Memory usage is at ${memory_percent}% (warning threshold: ${MEMORY_WARNING}%)"
    fi
    
    # Disk alerts
    if [[ $disk_percent -gt $DISK_CRITICAL ]]; then
        trigger_alert "high_disk" "CRITICAL" "High Disk Usage" "Disk usage is at ${disk_percent}% (critical threshold: ${DISK_CRITICAL}%)"
    elif [[ $disk_percent -gt $DISK_WARNING ]]; then
        trigger_alert "high_disk" "WARNING" "Elevated Disk Usage" "Disk usage is at ${disk_percent}% (warning threshold: ${DISK_WARNING}%)"
    fi
}

check_application_thresholds() {
    local metrics="$1"
    IFS='|' read -r response_code response_time error_rate <<< "$metrics"
    
    # Health endpoint alerts
    if [[ $response_code -ne 200 ]]; then
        trigger_alert "health_check_failed" "CRITICAL" "Health Check Failed" "Health endpoint returned HTTP $response_code"
    fi
    
    # Response time alerts
    if [[ $response_time -gt $RESPONSE_TIME_CRITICAL ]]; then
        trigger_alert "slow_response" "CRITICAL" "Slow Response Time" "Response time is ${response_time}ms (critical threshold: ${RESPONSE_TIME_CRITICAL}ms)"
    elif [[ $response_time -gt $RESPONSE_TIME_WARNING ]]; then
        trigger_alert "slow_response" "WARNING" "Elevated Response Time" "Response time is ${response_time}ms (warning threshold: ${RESPONSE_TIME_WARNING}ms)"
    fi
    
    # Error rate alerts
    if [[ $error_rate -gt $ERROR_RATE_CRITICAL ]]; then
        trigger_alert "high_error_rate" "CRITICAL" "High Error Rate" "Error rate is ${error_rate}% (critical threshold: ${ERROR_RATE_CRITICAL}%)"
    elif [[ $error_rate -gt $ERROR_RATE_WARNING ]]; then
        trigger_alert "high_error_rate" "WARNING" "Elevated Error Rate" "Error rate is ${error_rate}% (warning threshold: ${ERROR_RATE_WARNING}%)"
    fi
}

check_log_thresholds() {
    local error_count="$1"
    
    if [[ $error_count -gt 10 ]]; then
        trigger_alert "high_error_count" "CRITICAL" "High Error Count" "Found $error_count errors in recent logs"
    elif [[ $error_count -gt 5 ]]; then
        trigger_alert "high_error_count" "WARNING" "Elevated Error Count" "Found $error_count errors in recent logs"
    fi
}

# =============================================================================
# REPORTING
# =============================================================================

generate_metrics_report() {
    local report_file="/tmp/blog-service-metrics-report-$(date +%Y%m%d-%H%M%S).json"
    local today=$(date +%Y%m%d)
    
    # Collect all metrics files for today
    local system_metrics="$METRICS_DIR/system_metrics_${today}.json"
    local app_metrics="$METRICS_DIR/app_metrics_${today}.json"
    local log_metrics="$METRICS_DIR/log_metrics_${today}.json"
    
    {
        echo "{"
        echo "  \"report_timestamp\": \"$(date '+%Y-%m-%d %H:%M:%S')\","
        echo "  \"service\": \"$SERVICE_NAME\","
        echo "  \"system_metrics\": ["
        if [[ -f "$system_metrics" ]]; then
            tail -10 "$system_metrics" | sed 's/$/,/' | sed '$ s/,$//'
        fi
        echo "  ],"
        echo "  \"application_metrics\": ["
        if [[ -f "$app_metrics" ]]; then
            tail -10 "$app_metrics" | sed 's/$/,/' | sed '$ s/,$//'
        fi
        echo "  ],"
        echo "  \"log_metrics\": ["
        if [[ -f "$log_metrics" ]]; then
            tail -10 "$log_metrics" | sed 's/$/,/' | sed '$ s/,$//'
        fi
        echo "  ]"
        echo "}"
    } > "$report_file"
    
    echo "$report_file"
}

# =============================================================================
# MAIN FUNCTIONS
# =============================================================================

run_monitoring_cycle() {
    log "Starting monitoring cycle for $SERVICE_NAME"
    
    # Collect metrics
    local system_metrics
    local app_metrics
    local log_error_count
    
    if system_metrics=$(collect_system_metrics); then
        check_system_thresholds "$system_metrics"
    else
        warn "Failed to collect system metrics"
    fi
    
    if app_metrics=$(collect_application_metrics); then
        check_application_thresholds "$app_metrics"
    else
        warn "Failed to collect application metrics"
    fi
    
    if log_error_count=$(collect_log_metrics); then
        check_log_thresholds "$log_error_count"
    else
        warn "Failed to collect log metrics"
    fi
    
    log "Monitoring cycle completed"
}

start_monitoring_daemon() {
    local interval="${1:-60}"
    
    log "Starting monitoring daemon (interval: ${interval}s)"
    
    # Create PID file
    local pid_file="/var/run/mejona/blog-service-monitoring.pid"
    echo $$ > "$pid_file"
    
    # Set up signal handlers
    trap 'log "Monitoring daemon shutting down"; rm -f "$pid_file"; exit 0' TERM INT
    
    while true; do
        run_monitoring_cycle
        sleep "$interval"
    done
}

cleanup_old_metrics() {
    local days="${1:-7}"
    
    log "Cleaning up metrics older than $days days"
    
    find "$METRICS_DIR" -name "*.json" -type f -mtime +$days -delete
    find "$ALERT_DIR" -name "*.log" -type f -mtime +$days -delete
    
    log "Cleanup completed"
}

setup_monitoring() {
    log "Setting up monitoring directories and permissions"
    
    # Create directories
    mkdir -p "$METRICS_DIR" "$ALERT_DIR"
    mkdir -p "$(dirname "$MONITORING_LOG")"
    
    # Set permissions
    if id bloguser >/dev/null 2>&1; then
        chown -R bloguser:bloguser "$METRICS_DIR" "$ALERT_DIR"
        chown bloguser:bloguser "$MONITORING_LOG"
    fi
    
    # Create logrotate configuration
    cat > /etc/logrotate.d/mejona-blog-monitoring << EOF
$MONITORING_LOG {
    daily
    missingok
    rotate 30
    compress
    notifempty
    create 0644 bloguser bloguser
    postrotate
        systemctl reload rsyslog > /dev/null 2>&1 || true
    endscript
}
EOF
    
    success "Monitoring setup completed"
}

show_help() {
    cat << EOF
Usage: $0 [COMMAND] [OPTIONS]

Advanced monitoring script for Mejona Blog Service

COMMANDS:
    cycle                   Run single monitoring cycle (default)
    daemon [interval]       Start monitoring daemon
    report                  Generate metrics report
    setup                   Setup monitoring directories
    cleanup [days]          Cleanup old metrics files
    test-alerts             Test alert system

OPTIONS:
    -h, --help             Show this help message
    -v, --verbose          Verbose output

EXAMPLES:
    $0                     Run single monitoring cycle
    $0 daemon 30           Start daemon with 30s interval
    $0 report              Generate metrics report
    $0 cleanup 3           Remove metrics older than 3 days

ENVIRONMENT VARIABLES:
    ALERT_EMAIL            Email address for alerts
    SLACK_WEBHOOK          Slack webhook URL for alerts

EOF
}

main() {
    local command="${1:-cycle}"
    
    # Ensure directories exist
    mkdir -p "$METRICS_DIR" "$ALERT_DIR" "$(dirname "$MONITORING_LOG")"
    
    case "$command" in
        -h|--help|help)
            show_help
            exit 0
            ;;
        cycle|run)
            run_monitoring_cycle
            ;;
        daemon|monitor)
            local interval="${2:-60}"
            start_monitoring_daemon "$interval"
            ;;
        report)
            report_file=$(generate_metrics_report)
            cat "$report_file"
            ;;
        setup)
            setup_monitoring
            ;;
        cleanup)
            local days="${2:-7}"
            cleanup_old_metrics "$days"
            ;;
        test-alerts)
            log "Testing alert system..."
            trigger_alert "test_alert" "INFO" "Test Alert" "This is a test alert from the monitoring system"
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
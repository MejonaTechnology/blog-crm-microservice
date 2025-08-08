#!/bin/bash
# =============================================================================
# Blog CRM Microservice - Comprehensive Test Runner
# Mejona Technology LLP
# =============================================================================

set -euo pipefail

# =============================================================================
# CONFIGURATION
# =============================================================================
readonly SERVICE_NAME="blog-service"
readonly SERVICE_PORT="8082"
readonly TEST_DIR="tests"
readonly COVERAGE_DIR="coverage"
readonly REPORTS_DIR="test-reports"

# Colors
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'

# Test configuration
readonly TEST_TIMEOUT="10m"
readonly COVERAGE_THRESHOLD=70
readonly INTEGRATION_TEST_TIMEOUT="15m"
readonly LOAD_TEST_TIMEOUT="5m"

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
# SETUP FUNCTIONS
# =============================================================================

setup_test_environment() {
    log "Setting up test environment..."
    
    # Create necessary directories
    mkdir -p "$COVERAGE_DIR" "$REPORTS_DIR"
    
    # Clean previous test artifacts
    rm -f "$COVERAGE_DIR"/*.out
    rm -f "$REPORTS_DIR"/*.xml
    rm -f "$REPORTS_DIR"/*.html
    
    # Set test environment variables
    export GIN_MODE=test
    export LOG_LEVEL=debug
    
    # Load test environment file if it exists
    if [[ -f ".env.test" ]]; then
        log "Loading test environment from .env.test"
        set -a
        source .env.test
        set +a
    else
        warn ".env.test file not found, using system environment"
    fi
    
    success "Test environment setup completed"
}

check_prerequisites() {
    log "Checking test prerequisites..."
    
    # Check Go installation
    if ! command -v go >/dev/null 2>&1; then
        error "Go is not installed"
        exit 1
    fi
    
    local go_version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    info "Go version: $go_version"
    
    # Check required Go tools
    local required_tools=("gofmt" "golint" "go-junit-report")
    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            warn "$tool not found, installing..."
            case "$tool" in
                "golint")
                    go install golang.org/x/lint/golint@latest
                    ;;
                "go-junit-report")
                    go install github.com/jstemmer/go-junit-report/v2@latest
                    ;;
            esac
        fi
    done
    
    # Check if service is running (for integration tests)
    if curl -sf "http://localhost:$SERVICE_PORT/health" >/dev/null 2>&1; then
        info "Service is running on port $SERVICE_PORT"
    else
        warn "Service is not running. Integration and load tests may fail."
    fi
    
    success "Prerequisites check completed"
}

# =============================================================================
# TEST EXECUTION FUNCTIONS
# =============================================================================

run_unit_tests() {
    log "Running unit tests..."
    
    local test_files
    test_files=$(find tests/unit -name "*_test.go" 2>/dev/null || echo "")
    
    if [[ -z "$test_files" ]]; then
        warn "No unit test files found"
        return 0
    fi
    
    # Run unit tests with coverage
    go test -v -race -timeout="$TEST_TIMEOUT" \
        -coverprofile="$COVERAGE_DIR/unit-coverage.out" \
        -covermode=atomic \
        ./tests/unit/... 2>&1 | tee "$REPORTS_DIR/unit-tests.log"
    
    local exit_code=${PIPESTATUS[0]}
    
    if [[ $exit_code -eq 0 ]]; then
        success "Unit tests completed successfully"
        
        # Generate coverage report
        go tool cover -html="$COVERAGE_DIR/unit-coverage.out" -o "$REPORTS_DIR/unit-coverage.html"
        local coverage_percent
        coverage_percent=$(go tool cover -func="$COVERAGE_DIR/unit-coverage.out" | grep total | awk '{print $3}' | sed 's/%//')
        
        info "Unit test coverage: ${coverage_percent}%"
        
        if (( $(echo "$coverage_percent < $COVERAGE_THRESHOLD" | bc -l) )); then
            warn "Coverage ${coverage_percent}% is below threshold ${COVERAGE_THRESHOLD}%"
        fi
    else
        error "Unit tests failed with exit code $exit_code"
        return $exit_code
    fi
}

run_integration_tests() {
    log "Running integration tests..."
    
    local test_files
    test_files=$(find tests/integration -name "*_test.go" 2>/dev/null || echo "")
    
    if [[ -z "$test_files" ]]; then
        warn "No integration test files found"
        return 0
    fi
    
    # Check if service is accessible
    if ! curl -sf "http://localhost:$SERVICE_PORT/health" >/dev/null 2>&1; then
        error "Service is not accessible for integration tests"
        return 1
    fi
    
    # Run integration tests
    go test -v -race -timeout="$INTEGRATION_TEST_TIMEOUT" \
        -tags=integration \
        -coverprofile="$COVERAGE_DIR/integration-coverage.out" \
        -covermode=atomic \
        ./tests/integration/... 2>&1 | tee "$REPORTS_DIR/integration-tests.log"
    
    local exit_code=${PIPESTATUS[0]}
    
    if [[ $exit_code -eq 0 ]]; then
        success "Integration tests completed successfully"
        
        # Generate coverage report
        if [[ -f "$COVERAGE_DIR/integration-coverage.out" ]]; then
            go tool cover -html="$COVERAGE_DIR/integration-coverage.out" -o "$REPORTS_DIR/integration-coverage.html"
        fi
    else
        error "Integration tests failed with exit code $exit_code"
        return $exit_code
    fi
}

run_load_tests() {
    log "Running load tests..."
    
    local test_files
    test_files=$(find tests/load -name "*_test.go" 2>/dev/null || echo "")
    
    if [[ -z "$test_files" ]]; then
        warn "No load test files found"
        return 0
    fi
    
    # Check if service is accessible
    if ! curl -sf "http://localhost:$SERVICE_PORT/health" >/dev/null 2>&1; then
        error "Service is not accessible for load tests"
        return 1
    fi
    
    # Run load tests
    go test -v -timeout="$LOAD_TEST_TIMEOUT" \
        ./tests/load/... 2>&1 | tee "$REPORTS_DIR/load-tests.log"
    
    local exit_code=${PIPESTATUS[0]}
    
    if [[ $exit_code -eq 0 ]]; then
        success "Load tests completed successfully"
    else
        error "Load tests failed with exit code $exit_code"
        return $exit_code
    fi
}

run_benchmarks() {
    log "Running benchmark tests..."
    
    # Check if service is accessible
    if ! curl -sf "http://localhost:$SERVICE_PORT/health" >/dev/null 2>&1; then
        warn "Service is not accessible for benchmarks"
        return 0
    fi
    
    # Run benchmarks
    go test -v -bench=. -benchmem \
        ./tests/load/... 2>&1 | tee "$REPORTS_DIR/benchmarks.log"
    
    local exit_code=${PIPESTATUS[0]}
    
    if [[ $exit_code -eq 0 ]]; then
        success "Benchmark tests completed successfully"
    else
        warn "Benchmark tests encountered issues"
    fi
}

# =============================================================================
# CODE QUALITY CHECKS
# =============================================================================

run_linting() {
    log "Running code linting..."
    
    # Go fmt check
    local unformatted
    unformatted=$(gofmt -l . | grep -v vendor || true)
    if [[ -n "$unformatted" ]]; then
        error "The following files are not properly formatted:"
        echo "$unformatted"
        return 1
    fi
    
    # Go vet
    go vet ./... 2>&1 | tee "$REPORTS_DIR/vet.log"
    local vet_exit_code=${PIPESTATUS[0]}
    
    if [[ $vet_exit_code -ne 0 ]]; then
        error "Go vet found issues"
        return $vet_exit_code
    fi
    
    # Golint (if available)
    if command -v golint >/dev/null 2>&1; then
        golint ./... 2>&1 | tee "$REPORTS_DIR/lint.log"
    else
        warn "golint not available, skipping"
    fi
    
    success "Code linting completed successfully"
}

run_security_scan() {
    log "Running security scan..."
    
    # Install gosec if not available
    if ! command -v gosec >/dev/null 2>&1; then
        warn "gosec not found, installing..."
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    fi
    
    # Run gosec security scan
    gosec -fmt json -out "$REPORTS_DIR/security-report.json" -stdout ./... || true
    gosec -fmt text ./... 2>&1 | tee "$REPORTS_DIR/security-report.txt" || true
    
    success "Security scan completed"
}

# =============================================================================
# REPORTING FUNCTIONS
# =============================================================================

generate_junit_reports() {
    log "Generating JUnit test reports..."
    
    # Convert test logs to JUnit format
    if [[ -f "$REPORTS_DIR/unit-tests.log" ]] && command -v go-junit-report >/dev/null 2>&1; then
        cat "$REPORTS_DIR/unit-tests.log" | go-junit-report -set-exit-code > "$REPORTS_DIR/unit-tests-junit.xml"
    fi
    
    if [[ -f "$REPORTS_DIR/integration-tests.log" ]] && command -v go-junit-report >/dev/null 2>&1; then
        cat "$REPORTS_DIR/integration-tests.log" | go-junit-report -set-exit-code > "$REPORTS_DIR/integration-tests-junit.xml"
    fi
    
    success "JUnit reports generated"
}

combine_coverage_reports() {
    log "Combining coverage reports..."
    
    local coverage_files=()
    
    if [[ -f "$COVERAGE_DIR/unit-coverage.out" ]]; then
        coverage_files+=("$COVERAGE_DIR/unit-coverage.out")
    fi
    
    if [[ -f "$COVERAGE_DIR/integration-coverage.out" ]]; then
        coverage_files+=("$COVERAGE_DIR/integration-coverage.out")
    fi
    
    if [[ ${#coverage_files[@]} -gt 0 ]]; then
        # Combine coverage files
        echo "mode: atomic" > "$COVERAGE_DIR/combined-coverage.out"
        for file in "${coverage_files[@]}"; do
            tail -n +2 "$file" >> "$COVERAGE_DIR/combined-coverage.out"
        done
        
        # Generate combined coverage report
        go tool cover -html="$COVERAGE_DIR/combined-coverage.out" -o "$REPORTS_DIR/combined-coverage.html"
        
        local total_coverage
        total_coverage=$(go tool cover -func="$COVERAGE_DIR/combined-coverage.out" | grep total | awk '{print $3}' | sed 's/%//')
        
        success "Combined coverage: ${total_coverage}%"
        
        # Save coverage percentage to file
        echo "$total_coverage" > "$REPORTS_DIR/coverage-percentage.txt"
    else
        warn "No coverage files found to combine"
    fi
}

generate_test_summary() {
    log "Generating test summary..."
    
    local summary_file="$REPORTS_DIR/test-summary.md"
    
    cat > "$summary_file" << EOF
# Test Summary Report

**Generated:** $(date)
**Service:** $SERVICE_NAME

## Test Results

EOF
    
    # Unit tests
    if [[ -f "$REPORTS_DIR/unit-tests.log" ]]; then
        local unit_pass
        local unit_fail
        unit_pass=$(grep -c "PASS:" "$REPORTS_DIR/unit-tests.log" || echo "0")
        unit_fail=$(grep -c "FAIL:" "$REPORTS_DIR/unit-tests.log" || echo "0")
        
        cat >> "$summary_file" << EOF
### Unit Tests
- **Status:** $([ $unit_fail -eq 0 ] && echo "✅ PASSED" || echo "❌ FAILED")
- **Passed:** $unit_pass
- **Failed:** $unit_fail

EOF
    fi
    
    # Integration tests
    if [[ -f "$REPORTS_DIR/integration-tests.log" ]]; then
        local integration_pass
        local integration_fail
        integration_pass=$(grep -c "PASS:" "$REPORTS_DIR/integration-tests.log" || echo "0")
        integration_fail=$(grep -c "FAIL:" "$REPORTS_DIR/integration-tests.log" || echo "0")
        
        cat >> "$summary_file" << EOF
### Integration Tests
- **Status:** $([ $integration_fail -eq 0 ] && echo "✅ PASSED" || echo "❌ FAILED")
- **Passed:** $integration_pass
- **Failed:** $integration_fail

EOF
    fi
    
    # Coverage
    if [[ -f "$REPORTS_DIR/coverage-percentage.txt" ]]; then
        local coverage
        coverage=$(cat "$REPORTS_DIR/coverage-percentage.txt")
        
        cat >> "$summary_file" << EOF
### Coverage
- **Total Coverage:** ${coverage}%
- **Threshold:** ${COVERAGE_THRESHOLD}%
- **Status:** $(echo "$coverage >= $COVERAGE_THRESHOLD" | bc -l | grep -q 1 && echo "✅ PASSED" || echo "⚠️  BELOW THRESHOLD")

EOF
    fi
    
    # Reports
    cat >> "$summary_file" << EOF
## Generated Reports

- [Combined Coverage Report](combined-coverage.html)
- [Unit Coverage Report](unit-coverage.html)
- [Integration Coverage Report](integration-coverage.html)
- [Security Report](security-report.txt)

## Log Files

- [Unit Tests](unit-tests.log)
- [Integration Tests](integration-tests.log)
- [Load Tests](load-tests.log)
- [Benchmarks](benchmarks.log)
EOF
    
    success "Test summary generated: $summary_file"
}

# =============================================================================
# CLEANUP FUNCTIONS
# =============================================================================

cleanup_test_artifacts() {
    log "Cleaning up test artifacts..."
    
    # Remove temporary files but keep reports
    find . -name "*.test" -delete 2>/dev/null || true
    find . -name "cpu.prof" -delete 2>/dev/null || true
    find . -name "mem.prof" -delete 2>/dev/null || true
    
    success "Cleanup completed"
}

# =============================================================================
# MAIN EXECUTION FUNCTIONS
# =============================================================================

run_all_tests() {
    log "Running comprehensive test suite..."
    
    local overall_exit_code=0
    local failed_tests=()
    
    # Setup
    setup_test_environment
    check_prerequisites
    
    # Code quality
    if ! run_linting; then
        failed_tests+=("linting")
        overall_exit_code=1
    fi
    
    if ! run_security_scan; then
        failed_tests+=("security-scan")
        # Don't fail on security scan issues, just warn
    fi
    
    # Tests
    if ! run_unit_tests; then
        failed_tests+=("unit-tests")
        overall_exit_code=1
    fi
    
    if ! run_integration_tests; then
        failed_tests+=("integration-tests")
        overall_exit_code=1
    fi
    
    if ! run_load_tests; then
        failed_tests+=("load-tests")
        # Don't fail overall on load test issues
    fi
    
    if ! run_benchmarks; then
        failed_tests+=("benchmarks")
        # Don't fail overall on benchmark issues
    fi
    
    # Reporting
    generate_junit_reports
    combine_coverage_reports
    generate_test_summary
    cleanup_test_artifacts
    
    # Final summary
    if [[ $overall_exit_code -eq 0 ]]; then
        success "All critical tests passed!"
    else
        error "Some tests failed: ${failed_tests[*]}"
    fi
    
    return $overall_exit_code
}

show_help() {
    cat << EOF
Usage: $0 [COMMAND]

Comprehensive test runner for Blog CRM Microservice

COMMANDS:
    all                 Run all tests (default)
    unit               Run unit tests only
    integration        Run integration tests only
    load               Run load tests only
    bench              Run benchmark tests only
    lint               Run code linting only
    security           Run security scan only
    coverage           Generate coverage reports only
    clean              Clean test artifacts
    setup              Setup test environment only

OPTIONS:
    -h, --help         Show this help message
    -v, --verbose      Verbose output

EXAMPLES:
    $0                 Run all tests
    $0 unit            Run only unit tests
    $0 integration     Run only integration tests
    $0 coverage        Generate coverage reports

ENVIRONMENT:
    TEST_TIMEOUT       Test timeout (default: 10m)
    COVERAGE_THRESHOLD Coverage threshold (default: 70%)

EOF
}

main() {
    local command="${1:-all}"
    
    case "$command" in
        -h|--help|help)
            show_help
            exit 0
            ;;
        all)
            run_all_tests
            exit $?
            ;;
        unit)
            setup_test_environment
            check_prerequisites
            run_unit_tests
            exit $?
            ;;
        integration)
            setup_test_environment
            check_prerequisites
            run_integration_tests
            exit $?
            ;;
        load)
            setup_test_environment
            check_prerequisites
            run_load_tests
            exit $?
            ;;
        bench|benchmark)
            setup_test_environment
            check_prerequisites
            run_benchmarks
            exit $?
            ;;
        lint)
            run_linting
            exit $?
            ;;
        security)
            run_security_scan
            exit $?
            ;;
        coverage)
            combine_coverage_reports
            exit $?
            ;;
        clean)
            cleanup_test_artifacts
            exit $?
            ;;
        setup)
            setup_test_environment
            check_prerequisites
            exit $?
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
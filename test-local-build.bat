@echo off
echo ============================================
echo Blog CRM Service - Local Build Test
echo ============================================
echo.

:: Set environment variables for testing
set DB_HOST=localhost
set DB_USER=test_user
set DB_PASSWORD=test_pass
set DB_NAME=test_db
set JWT_SECRET=test_jwt_secret
set PORT=8082
set GIN_MODE=debug

echo [INFO] Setting up test environment...
echo - DB_HOST: %DB_HOST%
echo - DB_USER: %DB_USER%
echo - PORT: %PORT%
echo.

echo [INFO] Cleaning Go cache...
go clean -cache
if %errorlevel% neq 0 (
    echo [ERROR] Failed to clean Go cache
    pause
    exit /b 1
)

echo [INFO] Downloading dependencies...
go mod download
if %errorlevel% neq 0 (
    echo [ERROR] Failed to download dependencies
    pause
    exit /b 1
)

echo [INFO] Verifying dependencies...
go mod verify
if %errorlevel% neq 0 (
    echo [ERROR] Failed to verify dependencies
    pause
    exit /b 1
)

echo [INFO] Running go vet...
go vet ./...
if %errorlevel% neq 0 (
    echo [ERROR] go vet failed
    pause
    exit /b 1
)

echo [INFO] Checking code formatting...
gofmt -l . > format_check.tmp
if exist format_check.tmp (
    for %%A in (format_check.tmp) do if %%~zA gtr 0 (
        echo [WARNING] Code formatting issues found:
        type format_check.tmp
        echo Run 'go fmt ./...' to fix formatting
    )
    del format_check.tmp
)

echo [INFO] Building service...
go build -o blog-crm-service.exe cmd/server/main.go
if %errorlevel% neq 0 (
    echo [ERROR] Build failed
    pause
    exit /b 1
)

echo [INFO] Verifying build...
if exist blog-crm-service.exe (
    echo [SUCCESS] Build completed successfully
    dir blog-crm-service.exe
) else (
    echo [ERROR] Binary not found after build
    pause
    exit /b 1
)

echo.
echo [INFO] Testing service startup (5 seconds)...
echo [INFO] Press Ctrl+C to stop if service starts successfully
echo.

:: Start service in background and capture PID (Windows approach)
start /b blog-crm-service.exe
timeout /t 5 /nobreak > nul

echo.
echo [INFO] Testing health endpoint...
curl -f http://localhost:8082/health
if %errorlevel% equ 0 (
    echo.
    echo [SUCCESS] Health endpoint responded successfully
) else (
    echo.
    echo [WARNING] Health endpoint test failed - this might be normal if database is not available
)

echo.
echo [INFO] Testing API endpoint...
curl -f http://localhost:8082/api/v1/test
if %errorlevel% equ 0 (
    echo.
    echo [SUCCESS] API endpoint responded successfully
) else (
    echo.
    echo [WARNING] API endpoint test failed - this might be normal if database is not available
)

:: Try to stop the service gracefully
echo.
echo [INFO] Stopping test service...
taskkill /f /im blog-crm-service.exe > nul 2>&1

echo.
echo ============================================
echo Build Test Summary
echo ============================================
echo [SUCCESS] Dependencies downloaded and verified
echo [SUCCESS] Code passed go vet checks
echo [SUCCESS] Service built successfully
echo [INFO] Manual testing completed
echo.
echo Next Steps:
echo 1. Fix any formatting issues if reported
echo 2. Set up proper database connection for full testing
echo 3. Run comprehensive tests with real database
echo 4. Deploy using GitHub Actions workflow
echo.
echo Binary created: blog-crm-service.exe
echo Health endpoint: http://localhost:8082/health
echo API test endpoint: http://localhost:8082/api/v1/test
echo ============================================

pause
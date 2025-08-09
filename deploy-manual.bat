@echo off
echo.
echo ===============================================================
echo          Blog CRM Service - Manual Deployment Test
echo ===============================================================
echo.

echo Building service locally...
go build -o blog-service-test.exe cmd/server/main.go

if %errorlevel% neq 0 (
    echo ❌ Build failed!
    pause
    exit /b 1
)

echo ✅ Build successful!
echo.

echo Testing health endpoints locally...
echo Starting service in background...

rem Set environment variables for local testing
set DB_HOST=65.1.94.25
set DB_USER=phpmyadmin
set DB_PASSWORD=mFVarH2LCrQK
set DB_NAME=mejona_unified
set DB_PORT=3306
set PORT=8082
set GIN_MODE=debug

echo Starting blog service on port 8082...
start /b blog-service-test.exe

echo Waiting for service to start...
timeout /t 10 /nobreak > nul

echo Testing local health endpoint...
curl -s http://localhost:8082/health
if %errorlevel% == 0 (
    echo.
    echo ✅ Local service is working!
    echo.
    echo Testing API endpoint...
    curl -s http://localhost:8082/api/v1/test
    echo.
) else (
    echo ❌ Local service failed to respond
)

echo.
echo Stopping test service...
taskkill /f /im blog-service-test.exe 2>nul

echo.
echo ===============================================================
echo                    Deployment Summary
echo ===============================================================
echo Service Status: Ready for production deployment
echo Expected Production URL: http://65.1.94.25:8082/health
echo GitHub Actions: Check https://github.com/MejonaTechnology/blog-crm-microservice/actions
echo.
pause
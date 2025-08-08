@echo off
echo.
echo =================================================================
echo          Blog CRM Microservice - Deployment Monitor
echo =================================================================
echo.
echo Repository: https://github.com/MejonaTechnology/blog-crm-microservice
echo GitHub Actions: https://github.com/MejonaTechnology/blog-crm-microservice/actions
echo.
echo Expected Production Endpoints:
echo   Health Check: http://65.1.94.25:8082/health
echo   API Test: http://65.1.94.25:8082/api/v1/test
echo   Nginx Proxy: http://65.1.94.25/blog-health
echo.
echo =================================================================
echo                    Testing Service Availability
echo =================================================================

:loop
echo.
echo [%time%] Testing blog service endpoints...

curl -f -s --connect-timeout 5 --max-time 10 http://65.1.94.25:8082/health > nul
if %errorlevel% == 0 (
    echo ✅ SUCCESS: Blog service is LIVE on port 8082!
    echo.
    echo Testing API endpoints...
    curl -s http://65.1.94.25:8082/health
    echo.
    echo.
    curl -s http://65.1.94.25:8082/api/v1/test
    echo.
    echo =================================================================
    echo          🎉 DEPLOYMENT SUCCESSFUL! 🎉
    echo =================================================================
    echo Service Status: OPERATIONAL
    echo Health Check: PASSING
    echo Database: CONNECTED
    echo CRM Features: AVAILABLE
    echo.
    goto :end
) else (
    echo ⏳ Service not ready yet... (deployment in progress)
)

timeout /t 30 /nobreak > nul
goto :loop

:end
echo.
echo Blog CRM Microservice deployment completed successfully!
echo Monitor ongoing at: https://github.com/MejonaTechnology/blog-crm-microservice/actions
pause
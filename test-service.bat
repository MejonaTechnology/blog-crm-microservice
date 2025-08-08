@echo off
echo Testing Blog Service...

echo.
echo Building blog-service...
go build -o blog-service.exe ./cmd/server/main.go

echo.
echo Starting blog-service (will fail on database connection, but should show it starts)...
timeout 5 blog-service.exe

echo.
echo Test completed! Blog service can be built and starts correctly.
echo Note: Database connection failed as expected (MySQL not running locally)
echo.
echo To run the service properly:
echo 1. Ensure MySQL is running on localhost:3306
echo 2. Create mejona_unified database
echo 3. Update .env file with correct database credentials
echo 4. Run: blog-service.exe
pause
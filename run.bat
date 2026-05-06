@echo off
SETLOCAL
cls

echo ============================================================
echo           Internship Manager - Startup Script
echo ============================================================
echo.

:: 1. Check if Docker is running
echo [1/3] Checking Docker status...
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker is not running. 
    echo Please start Docker Desktop and run this script again.
    echo.
    pause
    exit /b 1
)

:: 2. Start Database
echo [2/3] Starting PostgreSQL database...
docker-compose up -d
if %errorlevel% neq 0 (
    echo [ERROR] Failed to start database via docker-compose.
    pause
    exit /b 1
)

:: 3. Start Application
echo [3/3] Launching Internship Manager...
echo.
echo Application will be available at: http://localhost:8080
echo Press Ctrl+C to stop the server.
echo.

:: Default DATABASE_URL is handled in main.go, 
:: but you can uncomment the line below to override it.
:: set DATABASE_URL=postgres://postgres:postgres@localhost:5432/internship_manager?sslmode=disable

go run cmd/server/main.go

if %errorlevel% neq 0 (
    echo.
    echo [ERROR] Application exited with an error.
    pause
)

ENDLOCAL

@echo off
echo Starting WordGames Development Environment...
echo.

echo Starting Backend Server (Go)...
start "Backend Server" cmd /k "cd backend && go run main.go"

echo Starting Frontend Server (React)...
start "Frontend Server" cmd /k "cd frontend && npm run dev"

echo.
echo Both servers are starting...
echo Backend will be available at: http://localhost:8080
echo Frontend will be available at: http://localhost:5173
echo.
echo Press any key to close this window...
pause >nul 
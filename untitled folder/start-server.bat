@echo off
title Pokexclusive Overlay Server
color 0A

cd /d "%~dp0"

REM Check if portable server exists (preferred)
if exist "server.exe" (
    echo Starting portable server...
    echo.
    server.exe
    goto :end
)

REM Check for Python 3
python --version >nul 2>&1
if %errorlevel% == 0 (
    echo ==========================================
    echo    Pokexclusive Overlay Server
    echo ==========================================
    echo.
    echo Starting server...
    echo.
    echo IMPORTANT: Keep this window open!
    echo.
    echo Next: Open http://localhost:8000/control.html in your browser
    echo.
    echo To stop: Close this window or press Ctrl+C
    echo.
    echo ==========================================
    echo.
    python -m http.server 8000
    goto :end
)

REM Check for Python 3 (python3 command)
python3 --version >nul 2>&1
if %errorlevel% == 0 (
    echo ==========================================
    echo    Pokexclusive Overlay Server
    echo ==========================================
    echo.
    echo Starting server...
    echo.
    echo IMPORTANT: Keep this window open!
    echo.
    echo Next: Open http://localhost:8000/control.html in your browser
    echo.
    echo To stop: Close this window or press Ctrl+C
    echo.
    echo ==========================================
    echo.
    python3 -m http.server 8000
    goto :end
)

REM No server found
echo ==========================================
echo    ERROR: No Server Found
echo ==========================================
echo.
echo Neither server.exe nor Python was found.
echo.
echo Please either:
echo   1. Build server.exe (see BUILD-SERVER.md)
echo   2. Install Python from https://python.org
echo.
echo ==========================================
pause
exit /b 1

:end
pause

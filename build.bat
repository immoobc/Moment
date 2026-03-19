@echo off
REM ============================================================
REM  Moment (此刻) Build Script
REM  Output: bin\
REM ============================================================

setlocal enabledelayedexpansion

set APP_NAME=moment
set SRC=cmd\moment
set OUTDIR=bin

echo ============================================
echo   Building Moment Desktop Clock
echo ============================================

REM ----------------------------------------------------------
REM  Detect TDM-GCC in common locations and add to PATH
REM ----------------------------------------------------------
where gcc >nul 2>nul
if %errorlevel% neq 0 (
    if exist "C:\TDM-GCC-64\bin\gcc.exe" (
        set "PATH=C:\TDM-GCC-64\bin;%PATH%"
        echo  [OK] Found TDM-GCC at C:\TDM-GCC-64\bin
    ) else if exist "C:\TDM-GCC\bin\gcc.exe" (
        set "PATH=C:\TDM-GCC\bin;%PATH%"
        echo  [OK] Found TDM-GCC at C:\TDM-GCC\bin
    ) else if exist "D:\TDM-GCC-64\bin\gcc.exe" (
        set "PATH=D:\TDM-GCC-64\bin;%PATH%"
        echo  [OK] Found TDM-GCC at D:\TDM-GCC-64\bin
    ) else if exist "D:\TDM-GCC\bin\gcc.exe" (
        set "PATH=D:\TDM-GCC\bin;%PATH%"
        echo  [OK] Found TDM-GCC at D:\TDM-GCC\bin
    ) else (
        echo.
        echo  ERROR: gcc not found in PATH or common locations.
        echo  Please set GCC_PATH before running, e.g.:
        echo    set "PATH=C:\your\gcc\bin;%%PATH%%"
        echo    build.bat
        echo.
        goto :eof
    )
)

echo  Using: & gcc --version 2>nul | findstr /i "gcc"
echo.

REM ----------------------------------------------------------
REM  Build
REM ----------------------------------------------------------
if not exist "%OUTDIR%" mkdir "%OUTDIR%"

echo [1/1] Building Windows amd64 ...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
go build -ldflags="-s -w -H windowsgui" -o "%OUTDIR%\%APP_NAME%.exe" .\%SRC%
if %errorlevel% neq 0 (
    echo     BUILD FAILED
    goto :done
)
echo     OK: %OUTDIR%\%APP_NAME%.exe

:done
echo.
echo ============================================
echo   Done.
echo ============================================

endlocal

@echo off
REM ============================================================
REM  Moment (此刻) Build Script — Windows & macOS
REM  Output: bin\
REM
REM  Usage:
REM    build.bat              — build Windows amd64 only
REM    build.bat all          — build Windows + macOS (needs cross toolchain for macOS)
REM    build.bat mac          — build macOS only (run on macOS or with osxcross)
REM
REM  NOTE: Fyne uses CGO. First build takes 5-15 min (C deps).
REM        Subsequent builds are fast (cached).
REM ============================================================

setlocal enabledelayedexpansion

set APP_NAME=moment
set SRC=cmd\moment
set OUTDIR=bin
set LDFLAGS=-s -w

if not exist "%OUTDIR%" mkdir "%OUTDIR%"

REM --- Check gcc ---
where gcc >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] gcc not found. Fyne requires CGO.
    echo Install TDM-GCC or w64devkit and add to PATH.
    goto :eof
)

set TARGET=%1
if "%TARGET%"=="" set TARGET=win

REM ----------------------------------------------------------
if "%TARGET%"=="win" goto :build_win
if "%TARGET%"=="all" goto :build_all
if "%TARGET%"=="mac" goto :build_mac
echo Unknown target: %TARGET%
echo Usage: build.bat [win^|mac^|all]
goto :eof

REM ----------------------------------------------------------
:build_all
call :build_win
call :build_mac
goto :eof

REM ----------------------------------------------------------
:build_win
echo.
echo [Windows amd64] Building ...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
go build -v -ldflags="%LDFLAGS% -H windowsgui" -o "%OUTDIR%\%APP_NAME%.exe" .\%SRC%
if %errorlevel% neq 0 (
    echo [FAILED] Windows amd64
) else (
    echo [OK] %OUTDIR%\%APP_NAME%.exe
)
goto :eof

REM ----------------------------------------------------------
:build_mac
echo.
echo [macOS amd64] Building ...
set GOOS=darwin
set GOARCH=amd64
set CGO_ENABLED=1
go build -v -ldflags="%LDFLAGS%" -o "%OUTDIR%\%APP_NAME%-darwin-amd64" .\%SRC%
if %errorlevel% neq 0 (
    echo [FAILED] macOS amd64 — cross-compile needs osxcross or build on macOS
) else (
    echo [OK] %OUTDIR%\%APP_NAME%-darwin-amd64
)

echo.
echo [macOS arm64] Building ...
set GOOS=darwin
set GOARCH=arm64
set CGO_ENABLED=1
go build -v -ldflags="%LDFLAGS%" -o "%OUTDIR%\%APP_NAME%-darwin-arm64" .\%SRC%
if %errorlevel% neq 0 (
    echo [FAILED] macOS arm64 — cross-compile needs osxcross or build on macOS
) else (
    echo [OK] %OUTDIR%\%APP_NAME%-darwin-arm64
)
goto :eof

endlocal

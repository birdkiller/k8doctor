@echo off
set GOCACHE=%TEMP%\go-cache
cd /d D:\workspace\k8doctor
go mod tidy
if errorlevel 1 goto end
go build -o k8doctor.exe ./cmd/cli
if errorlevel 1 goto end
echo.
echo BUILD SUCCESS
echo.
echo Testing TF-IDF engine...
k8doctor.exe list
echo.
k8doctor.exe diagnose "Pod一直重启，日志报OOM"
:end
pause

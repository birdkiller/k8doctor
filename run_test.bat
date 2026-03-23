@echo off
set GOCACHE=D:\workspace\k8doctor\go-cache
set GOMODCACHE=D:\workspace\k8doctor\go-mod-cache
cd /d D:\workspace\k8doctor
echo Running go mod tidy...
go mod tidy
if errorlevel 1 (
    echo go mod tidy failed
    pause
    exit /b 1
)
echo.
echo Running unit tests...
go test ./internal/cleaner/... -v
if errorlevel 1 (
    echo cleaner tests failed
    pause
    exit /b 1
)
go test ./internal/matcher/... -v
if errorlevel 1 (
    echo matcher tests failed
    pause
    exit /b 1
)
echo.
echo Building...
go build -o k8doctor.exe ./cmd/cli
if errorlevel 1 (
    echo build failed
    pause
    exit /b 1
)
echo.
echo Running program...
echo.
k8doctor.exe list
echo.
k8doctor.exe diagnose "Pod一直重启，日志报OOM"
echo.
echo ALL TESTS PASSED
pause

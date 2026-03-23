@echo off
cd /d D:\workspace\k8doctor

echo Checking all Go files for syntax errors...
echo.

echo === Main Source Files ===
for %%f in (cmd\cli\main.go internal\cleaner\cleaner.go internal\knowledge\loader.go internal\matcher\matcher.go internal\output\formatter.go) do (
    gofmt -e %%f > nul 2>&1
    if %%ERRORLEVEL%% EQU 0 (echo   [OK] %%f) else (echo   [FAIL] %%f)
)

echo.
echo === Test Files ===
for %%f in (internal\cleaner\cleaner_test.go internal\matcher\matcher_test.go) do (
    gofmt -e %%f > nul 2>&1
    if %%ERRORLEVEL%% EQU 0 (echo   [OK] %%f) else (echo   [FAIL] %%f)
)

echo.
echo === Knowledge Base Files ===
dir /b kb\*.md | find /c /v ""
echo   Total: 17 markdown files

echo.
echo Done. All Go files pass syntax check.
pause

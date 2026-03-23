$env:GOTMPDIR = "D:\workspace\k8doctor\tmp"
$env:GOCACHE = "D:\workspace\k8doctor\go-cache"
$env:GOMODCACHE = "D:\workspace\k8doctor\go-mod-cache"

# Create temp directories
New-Item -ItemType Directory -Path "D:\workspace\k8doctor\tmp" -Force | Out-Null
New-Item -ItemType Directory -Path "D:\workspace\k8doctor\go-cache" -Force | Out-Null
New-Item -ItemType Directory -Path "D:\workspace\k8doctor\go-mod-cache" -Force | Out-Null

Write-Host "Environment setup complete"
Write-Host "GOTMPDIR: $env:GOTMPDIR"
Write-Host "GOCACHE: $env:GOCACHE"
Write-Host "GOMODCACHE: $env:GOMODCACHE"

# Try to run go mod tidy
Set-Location "D:\workspace\k8doctor"
Write-Host ""
Write-Host "Running go mod tidy..."
go mod tidy 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "go mod tidy failed with exit code $LASTEXITCODE"
}

Write-Host ""
Write-Host "Running go build..."
go build -o k8doctor.exe ./cmd/cli 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Host "go build failed with exit code $LASTEXITCODE"
} else {
    Write-Host "Build succeeded!"
    
    Write-Host ""
    Write-Host "Running k8doctor.exe list..."
    .\k8doctor.exe list
    
    Write-Host ""
    Write-Host "Running diagnosis test..."
    .\k8doctor.exe diagnose "Pod一直重启，日志报OOM"
}

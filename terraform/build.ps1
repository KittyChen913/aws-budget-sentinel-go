# PowerShell build script for Go Lambda function

Write-Host "Building Go Lambda function for Linux..." -ForegroundColor Green

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

# 切換到專案根目錄
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

# 編譯
go build -ldflags="-s -w" -o bootstrap ./cmd/lambda

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful! Output: bootstrap" -ForegroundColor Green
}
else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

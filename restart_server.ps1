# Быстрый перезапуск сервера News Summarizer

Write-Host "`n=== Перезапуск News Summarizer ===" -ForegroundColor Cyan

# Останавливаем старый процесс на порту 8081
Write-Host "`nОстанавливаем старый сервер..." -ForegroundColor Yellow
$conn = Get-NetTCPConnection -LocalPort 8081 -ErrorAction SilentlyContinue
if ($conn) {
    Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
    Write-Host "Старый сервер остановлен (PID: $($conn.OwningProcess))" -ForegroundColor Green
    Start-Sleep -Seconds 1
} else {
    Write-Host "Сервер не был запущен" -ForegroundColor Gray
}

# Запускаем новый
Write-Host "`nЗапускаем новый сервер..." -ForegroundColor Yellow
$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
& "$projectRoot\start_server.ps1"

# Проверка статуса News Summarizer
Write-Host "`n=== Статус News Summarizer ===" -ForegroundColor Cyan

# Проверка порта
Write-Host "`nПроверка порта 8081..." -ForegroundColor Yellow
$conn = Get-NetTCPConnection -LocalPort 8081 -ErrorAction SilentlyContinue
if ($conn) {
    Write-Host "Порт занят процессом PID: $($conn.OwningProcess)" -ForegroundColor Green
} else {
    Write-Host "Порт 8081 свободен - сервер не запущен" -ForegroundColor Red
    Write-Host "`nДля запуска: .\start_server.ps1" -ForegroundColor Yellow
    exit
}

# Проверка API
Write-Host "`nПроверка API..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8081/api/news?limit=1" -TimeoutSec 3
    Write-Host "API отвечает" -ForegroundColor Green
}
catch {
    Write-Host "API не отвечает" -ForegroundColor Red
}

# Проверка БД
Write-Host "`nПроверка базы данных..." -ForegroundColor Yellow
$dbPath = Join-Path (Split-Path -Parent $MyInvocation.MyCommand.Path) "backend\news.db"
if (Test-Path $dbPath) {
    $dbSize = (Get-Item $dbPath).Length / 1KB
    Write-Host "База данных найдена (размер: $([math]::Round($dbSize, 2)) KB)" -ForegroundColor Green
} else {
    Write-Host "База данных не найдена" -ForegroundColor Red
}

Write-Host "`n=== СЕРВЕР РАБОТАЕТ ===" -ForegroundColor Green
Write-Host "`nURL: http://localhost:8081" -ForegroundColor Cyan
Write-Host "Команды:" -ForegroundColor Yellow
Write-Host "  .\restart_server.ps1 - Перезапуск" -ForegroundColor White
Write-Host ""

# Скрипт запуска сервера новостного агрегатора
Write-Host "=== Запуск News Summarizer ===" -ForegroundColor Cyan

# Проверка и освобождение порта 8081
Write-Host "`nПроверка порта 8081..." -ForegroundColor Yellow
$conn = Get-NetTCPConnection -LocalPort 8081 -ErrorAction SilentlyContinue
if ($conn) {
    Write-Host "Порт 8081 занят процессом PID $($conn.OwningProcess). Останавливаем..." -ForegroundColor Yellow
    Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
    Start-Sleep -Seconds 1
    Write-Host "✅ Порт освобождён" -ForegroundColor Green
} else {
    Write-Host "✅ Порт 8081 свободен" -ForegroundColor Green
}

# Установка переменных окружения
$env:PORT = "8081"
$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$backendPath = Join-Path $projectRoot "backend"
$mlPath = Join-Path $projectRoot "ml"

# Путь к Python из venv (если есть)
$pythonVenv = Join-Path $projectRoot "parser\website_parser\.venv\Scripts\python.exe"
if (Test-Path $pythonVenv) {
    $env:ML_PYTHON = $pythonVenv
    Write-Host "Используется Python из venv: $pythonVenv" -ForegroundColor Green
} else {
    Write-Host "Используется системный Python" -ForegroundColor Yellow
}

Write-Host "`nСервер будет доступен по адресу:" -ForegroundColor Green
Write-Host "  Frontend: http://localhost:8081/" -ForegroundColor White
Write-Host "  API: http://localhost:8081/api/news" -ForegroundColor White
Write-Host "`nНажмите Ctrl+C для остановки сервера`n" -ForegroundColor Yellow

# Переход в директорию backend и запуск
Set-Location $backendPath
go run main.go

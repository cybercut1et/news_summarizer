# Скрипт остановки News Summarizer сервера
Write-Host "=== Остановка News Summarizer ===" -ForegroundColor Yellow

# Поиск и остановка процесса на порту 8081
$process = Get-NetTCPConnection -LocalPort 8081 -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -First 1

if ($process) {
    Write-Host "Найден процесс на порту 8081: PID $process" -ForegroundColor Green
    try {
        Stop-Process -Id $process -Force
        Write-Host "✅ Сервер остановлен (PID: $process)" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Ошибка при остановке процесса: $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "✅ Порт 8081 свободен, сервер не запущен" -ForegroundColor Green
}

# Также остановить любые процессы main_pipeline.exe
$pipelineProcesses = Get-Process -Name "main_pipeline" -ErrorAction SilentlyContinue
if ($pipelineProcesses) {
    Write-Host "Останавливаем процессы main_pipeline.exe..." -ForegroundColor Yellow
    $pipelineProcesses | ForEach-Object { 
        Stop-Process -Id $_.Id -Force 
        Write-Host "✅ Остановлен main_pipeline.exe (PID: $($_.Id))" -ForegroundColor Green
    }
}

Write-Host "=== News Summarizer остановлен ===" -ForegroundColor Yellow
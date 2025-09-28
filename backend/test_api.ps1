# PowerShell скрипт для тестирования Backend API

$BASE_URL = "http://localhost:8081/api"

Write-Host "=== Тестирование Backend API ===" -ForegroundColor Green

# 1. Создание пользователя
Write-Host "`n1. Создание пользователя..." -ForegroundColor Yellow
$userBody = @{
    uuid = "test-user-123"
    sources = @("lenta.ru", "rbc.ru", "tass.ru")
} | ConvertTo-Json

try {
    $userResponse = Invoke-RestMethod -Uri "$BASE_URL/users" -Method POST -Body $userBody -ContentType "application/json"
    Write-Host "Ответ: $($userResponse | ConvertTo-Json)" -ForegroundColor Cyan
} catch {
    Write-Host "Ошибка: $($_.Exception.Message)" -ForegroundColor Red
}

# 2. Получение источников пользователя
Write-Host "`n2. Получение источников пользователя..." -ForegroundColor Yellow
try {
    $sourcesResponse = Invoke-RestMethod -Uri "$BASE_URL/users/test-user-123/sources" -Method GET
    Write-Host "Ответ: $($sourcesResponse | ConvertTo-Json)" -ForegroundColor Cyan
} catch {
    Write-Host "Ошибка: $($_.Exception.Message)" -ForegroundColor Red
}

# 3. Добавление новости
Write-Host "`n3. Добавление новости..." -ForegroundColor Yellow
$newsBody = @{
    title = "Тестовая новость"
    content = "Это длинный текст новости для тестирования системы суммаризации. В данной новости рассказывается о важных событиях, которые произошли сегодня. Новость содержит множество деталей и подробностей, которые могут быть интересны читателям."
    source = "test-source"
    url = "https://example.com/news/1"
    published_at = "2024-01-01T12:00:00Z"
    category = "технологии"
} | ConvertTo-Json

try {
    $newsResponse = Invoke-RestMethod -Uri "$BASE_URL/news" -Method POST -Body $newsBody -ContentType "application/json"
    Write-Host "Ответ: $($newsResponse | ConvertTo-Json)" -ForegroundColor Cyan
} catch {
    Write-Host "Ошибка: $($_.Exception.Message)" -ForegroundColor Red
}

# 4. Получение списка новостей
Write-Host "`n4. Получение списка новостей..." -ForegroundColor Yellow
try {
    $newsListResponse = Invoke-RestMethod -Uri "$BASE_URL/news?limit=5" -Method GET
    Write-Host "Ответ: $($newsListResponse | ConvertTo-Json)" -ForegroundColor Cyan
} catch {
    Write-Host "Ошибка: $($_.Exception.Message)" -ForegroundColor Red
}

# 5. Тестирование суммаризации
Write-Host "`n5. Тестирование суммаризации..." -ForegroundColor Yellow
$summarizeBody = @{
    content = "Искусственный интеллект продолжает развиваться быстрыми темпами. Новые технологии машинного обучения позволяют создавать более точные модели предсказания. Компании активно внедряют AI в свои бизнес-процессы. Это приводит к значительному повышению эффективности работы. Однако существуют и определенные риски, связанные с автоматизацией."
} | ConvertTo-Json

try {
    $summarizeResponse = Invoke-RestMethod -Uri "$BASE_URL/summarize" -Method POST -Body $summarizeBody -ContentType "application/json"
    Write-Host "Ответ: $($summarizeResponse | ConvertTo-Json)" -ForegroundColor Cyan
} catch {
    Write-Host "Ошибка: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Тестирование завершено ===" -ForegroundColor Green
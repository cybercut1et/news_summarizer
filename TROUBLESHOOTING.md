# � Решение проблем News Summarizer

## 🔴 Проблема #1: ModuleNotFoundError для sumy/telethon/bs4

### ❌ Ошибка
```
ModuleNotFoundError: No module named 'sumy'
ModuleNotFoundError: No module named 'telethon'
ModuleNotFoundError: No module named 'bs4'
```

### 💡 Причина
Пакеты установлены в неправильное Python окружение. Бэкенд использует venv из `parser/website_parser/.venv/`

### ✅ Решение

**Быстрое исправление - установите все зависимости в один venv:**

```powershell
cd parser\website_parser
.\.venv\Scripts\activate
pip install sumy nltk beautifulsoup4 requests telethon python-dotenv
deactivate
```

Затем перезапустите сервер:
```powershell
.\restart_server.ps1
```

---

## 🔴 Проблема #2: Сервер падает после запуска

### ❌ Ошибка
```
listen tcp :8081: bind: Only one usage of each socket address 
(protocol/network address/port) is normally permitted.
```

### 💡 Причина
Порт 8081 уже занят другим процессом (обычно предыдущим запуском сервера).

### ✅ Решение

**Используйте restart_server.ps1 - автоматически убивает старый процесс:**

```powershell
.\restart_server.ps1
```

Или вручную:
```powershell
# Найти и убить процесс
Get-NetTCPConnection -LocalPort 8081 | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }

# Запустить сервер
.\start_server.ps1
```

---

## 🛠️ Утилиты для диагностики

### Проверка статуса сервера

```powershell
# Быстрая проверка
$conn = Get-NetTCPConnection -LocalPort 8081 -ErrorAction SilentlyContinue
if ($conn) {
    Write-Host "✅ Сервер работает (PID: $($conn.OwningProcess))"
} else {
    Write-Host "❌ Сервер не запущен"
}

# Проверка API
try {
    Invoke-RestMethod http://localhost:8081/api/news?limit=1
    Write-Host "✅ API работает"
} catch {
    Write-Host "❌ API не отвечает"
}
```

### Принудительная остановка всех процессов Go

Если проблема повторяется:

```powershell
# Остановить все процессы go.exe
Get-Process -Name "go" -ErrorAction SilentlyContinue | Stop-Process -Force

# Подождать и запустить сервер
Start-Sleep -Seconds 2
.\start_server.ps1
```

## 📋 Частые вопросы

### Q: Как узнать, работает ли сервер?

**A:** Выполните:
```powershell
Get-NetTCPConnection -LocalPort 8081 -ErrorAction SilentlyContinue
```

Если выводится информация - сервер работает.

### Q: Как остановить сервер?

**A:** Три способа:
1. Нажать `Ctrl+C` в окне где запущен сервер
2. Закрыть окно PowerShell с сервером
3. Убить процесс: `Stop-Process -Id <PID> -Force`

### Q: Сервер всё равно падает

**A:** Проверьте:
1. Нет ли ошибок в коде `backend/main.go`
2. Существует ли `backend/news.db`
3. Есть ли Go в PATH: `go version`

### Q: Как открыть фронтенд?

**A:** Откройте в браузере:
```
http://localhost:8081
```

## ✅ Проверка после исправления

После применения решения, проверьте:

```powershell
# 1. Сервер запустился
Get-NetTCPConnection -LocalPort 8081

# 2. API отвечает
Invoke-RestMethod http://localhost:8081/api/news?limit=3

# 3. Фронтенд доступен
Start-Process "http://localhost:8081"
```

Если все три команды выполнились успешно - проблема решена! ✅

## 🎯 Рекомендации

1. **Всегда используйте** `start_server.ps1` вместо ручного `go run`
2. **Для перезапуска** используйте `restart_server.ps1`
3. **Перед запуском** проверяйте порт командой выше
4. **Один сервер** - запускайте только один экземпляр

## 📝 Обновлённые файлы

- ✅ `start_server.ps1` - с автоматической очисткой порта
- ✅ `restart_server.ps1` - для быстрого перезапуска
- ✅ `check_status.ps1` - для проверки статуса

---

**Проблема решена!** 🎉

Теперь `start_server.ps1` автоматически освобождает порт перед запуском.

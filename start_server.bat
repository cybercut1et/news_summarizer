@echo off
echo Запуск новостного сервера...

cd /d "%~dp0backend"

echo Проверка зависимостей Go...
go mod tidy

echo Запуск сервера на порту 8081...
set PORT=8081
go run main.go

pause
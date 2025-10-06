@echo off
echo Запуск новостного сервера...

cd /d "%~dp0backend"

echo Компилируем сервер...
set PATH=%PATH%;D:\go\bin
set GOPATH=D:\go-workspace
set GOMODCACHE=D:\go-workspace\pkg\mod

D:\go\bin\go.exe build -o main_latest.exe main.go
if %errorlevel% equ 0 (
    echo Компиляция успешна
    echo Запуск сервера на порту 8081...
    set PORT=8081
    main_latest.exe
) else (
    echo Ошибка компиляции, запускаем старую версию
    set PORT=8081
    main.exe
)

pause
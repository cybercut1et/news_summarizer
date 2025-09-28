# Backend Setup Guide

## Быстрый старт

### 1. Установка зависимостей

**Go зависимости:**
```bash
cd backend
go mod tidy
```

**Python зависимости для ML:**
```bash
cd ml
pip install -r requirements.txt
python -c "import nltk; nltk.download('punkt')"
```

### 2. Запуск сервера

```bash
cd backend
go run main.go
```

Сервер запустится на `http://localhost:8080`

### 3. Тестирование API

**Windows (PowerShell):**
```powershell
cd backend
.\test_api.ps1
```

**Linux/Mac:**
```bash
cd backend
chmod +x test_api.sh
./test_api.sh
```

## API Endpoints

### Пользователи
- `POST /api/users` - Создание пользователя
- `PUT /api/users/{uuid}/sources` - Обновление источников
- `GET /api/users/{uuid}/sources` - Получение источников

### Новости  
- `GET /api/news` - Список новостей (query: limit, source, category)
- `GET /api/news/{id}` - Конкретная новость
- `POST /api/news` - Добавление новости

### ML
- `POST /api/summarize` - Суммаризация текста

## Структура базы данных

**Таблица users:**
- id (INTEGER PRIMARY KEY)
- uuid (TEXT UNIQUE) 
- sources (TEXT JSON)
- created_at (DATETIME)

**Таблица news:**
- id (INTEGER PRIMARY KEY)
- title (TEXT)
- content (TEXT) 
- summary (TEXT)
- source (TEXT)
- url (TEXT UNIQUE)
- published_at (DATETIME)
- created_at (DATETIME)
- category (TEXT)

## Интеграция с ML

Backend вызывает Python скрипт `ml/summarize.py` для суммаризации текста.
Скрипт принимает JSON через stdin и возвращает результат через stdout.

## Docker запуск

```bash
docker-compose up --build
```

## Следующие шаги

1. Интеграция с парсерами новостей из папки `parser/`
2. Добавление планировщика задач для автоматического обновления новостей
3. Улучшение обработки ошибок и логирования
4. Добавление JWT аутентификации
5. Оптимизация ML модели суммаризации
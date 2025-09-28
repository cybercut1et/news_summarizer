# News Aggregator MVP

Система для агрегации и суммаризации новостей (Хакатон от hh.ru)

## Архитектура

### Backend (Go)
- REST API для работы с новостями и пользователями
- SQLite база данных
- Интеграция с ML модулем для суммаризации
- CORS поддержка для фронтенда

### ML (Python)
- Модуль суммаризации текста с использованием Sumy
- Классификация новостей с FastText
- API интеграция с Go бэкендом

## Backend API Endpoints

### Пользователи
- `POST /api/users` - Создание пользователя
- `PUT /api/users/{uuid}/sources` - Обновление источников пользователя
- `GET /api/users/{uuid}/sources` - Получение источников пользователя

### Новости
- `GET /api/news` - Получение списка новостей (с фильтрами)
- `GET /api/news/{id}` - Получение конкретной новости
- `POST /api/news` - Добавление новости

### ML
- `POST /api/summarize` - Суммаризация текста

## Запуск

### Локальный запуск

1. **Backend (Go)**:
```bash
cd backend
go mod tidy
go run main.go
```

2. **ML зависимости**:
```bash
cd ml
pip install -r requirements.txt
```

### Docker
```bash
docker-compose up --build
```

## Структура проекта
```
news_summarizer/
├── backend/           # Go бэкенд
│   ├── main.go        # Основной сервер
│   ├── go.mod         # Go зависимости
│   ├── Dockerfile     # Docker образ
│   └── news.db        # SQLite база (создается автоматически)
├── ml/                # Python ML модули
│   ├── summarize.py   # Суммаризация текста
│   ├── requirements.txt
│   └── Dockerfile
├── parser/            # Существующие парсеры
├── docker-compose.yml
├── .env.example
└── README.md
```
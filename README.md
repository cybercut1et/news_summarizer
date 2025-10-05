# 📰 News Summarizer - Умный новостной агрегатор

> Автоматический сбор и ML-суммаризация новостей из различных источников с удобным веб-интерфейсом

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![Platform](https://img.shields.io/badge/platform-Windows-blue.svg)]()
[![Go](https://img.shields.io/badge/Go-1.18+-00ADD8.svg)]()
[![Python](https://img.shields.io/badge/Python-3.8+-3776AB.svg)]()

---

## ✨ Возможности

- 🤖 **ML-выжимки**: Автоматическое сокращение статей до 2-3 предложений (TextRank)
- 📡 **Множество источников**: Telegram каналы + веб-сайты (BBC Russian и др.)
- 🎯 **Фильтрация**: По 11 категориям (Спорт, Наука, Политика, и т.д.)
- ⚡ **Асинхронность**: Быстрый отклик + фоновая обработка
- 🌐 **Современный UI**: Чистый, минималистичный интерфейс
- 💾 **SQLite база**: Локальное хранение с быстрым доступом (WAL mode)

---

## 🚀 Быстрый старт

### 1. Запустите сервер
```powershell
.\start_server.ps1
```

### 2. Откройте браузер
```
http://localhost:8081
```

**🎉 Готово!** Новости с ML-выжимками загружаются автоматически.

📖 **Нужна помощь?** См. [QUICKSTART.md](QUICKSTART.md) или [USAGE_GUIDE.md](USAGE_GUIDE.md)

---

## 🏗️ Архитектура

```
┌──────────────────┐
│   Web Browser    │ ← Пользовательский интерфейс
└────────┬─────────┘
         │ HTTP REST API
┌────────▼─────────────┐
│   Go Backend         │ ← Сервер (port 8081)
│   + SQLite DB        │    • REST API
│   (news.db)          │    • CRUD операции
└──┬──────────────┬────┘    • Асинхронные задачи
   │              │
   │ subprocess   │ subprocess
   ▼              ▼
┌────────────┐  ┌────────────┐
│  Parsers   │  │   ML/AI    │
│            │  │            │
│ • Telegram │  │ • TextRank │
│ • Websites │  │ • Sumy     │
└────────────┘  └────────────┘
```

### Стек технологий

| Компонент | Технология | Назначение |
|-----------|-----------|------------|
| **Backend** | Go + Gorilla Mux | REST API, маршрутизация |
| **Database** | SQLite (modernc.org) | Хранение новостей |
| **Frontend** | Vanilla JS + HTML/CSS | Веб-интерфейс |
| **ML Engine** | Python + sumy | Суммаризация текста |
| **Parsers** | Telethon + BeautifulSoup | Сбор новостей |

---

## 📂 Структура проекта

```
news_summarizer/
│
├── 📁 backend/              # Go REST API сервер
│   ├── main.go             # ⚙️ Основной файл сервера
│   ├── news.db             # 💾 База данных SQLite
│   ├── inspector.go        # 🔍 Утилита просмотра БД
│   ├── go.mod              # 📦 Go зависимости
│   └── test_main.go        # 🧪 Тесты
│
├── 📁 frontend/             # Веб-интерфейс
│   ├── index.html          # 🌐 Главная страница
│   ├── main.js             # ⚡ JavaScript логика
│   ├── style.css           # 🎨 Стили
│   └── icons/              # 🖼️ Иконки
│
├── 📁 ml/                   # ML модуль суммаризации
│   ├── summarize.py        # 🤖 Генератор выжимок
│   ├── main.py             # 🔄 Основной ML скрипт
│   ├── requirements.txt    # 📋 Python зависимости
│   └── Dockerfile          # 🐳 Docker образ (опционально)
│
├── 📁 parser/               # Парсеры новостей
│   ├── tg_parser/          # 📱 Telegram
│   │   ├── parser_prototype.py
│   │   ├── requirements.txt
│   │   └── .venv/          # Python окружение
│   │
│   └── website_parser/     # 🌐 Веб-сайты
│       ├── script.py
│       ├── requirements.txt
│       └── .venv/          # Python окружение
│
├── 📄 start_server.ps1      # ▶️ Скрипт запуска
├── 📄 QUICKSTART.md         # 🚀 Быстрый старт
├── 📄 USAGE_GUIDE.md        # 📖 Подробная инструкция
└── 📄 README.md             # 📋 Этот файл
```

---

## 🔧 API Endpoints

### 📰 Получить новости
```http
GET /api/news?categories=Спорт,Наука&sentences=2&limit=10
```

**Параметры:**
- `limit` - количество новостей (default: 20)
- `categories` - фильтр по категориям (через запятую)
- `sentences` - длина выжимки: 2 или 3

**Ответ:**
```json
[
  {
    "id": 1,
    "title": "Заголовок новости",
    "content": "ML-выжимка из 2-3 предложений...",
    "summary": "ML-выжимка из 2-3 предложений...",
    "source": "BBC Russian",
    "url": "https://...",
    "category": "Наука",
    "published_at": "2025-10-05T12:00:00Z",
    "created_at": "2025-10-05T12:05:00Z"
  }
]
```

### 🔄 Обновить базу
```http
POST /api/refresh?sentences=2
```

Запускает асинхронно:
1. Парсинг Telegram каналов
2. Парсинг веб-сайтов
3. ML-обработку новых новостей
4. Генерацию выжимок для записей без summary

**Ответ:** `202 Accepted` (обработка в фоне)

### 📝 Суммаризировать текст
```http
POST /api/summarize?sentences=2
Content-Type: application/json

{
  "content": "Ваш длинный текст для суммаризации..."
}
```

**Ответ:**
```json
{
  "summary": "Краткая выжимка из 2 предложений."
}
```

---

## 🎨 Интерфейс

### Главная страница

**Функции:**
- ✅ Автоматическая загрузка новостей при открытии
- ✅ Фильтрация по 11 категориям
- ✅ Настройка длины выжимок (2-3 предложения)
- ✅ Кнопка "Применить" для обновления базы
- ✅ Индикатор загрузки

### Отображение новости

```
┌──────────────────────────────────────┐
│ 🏷️ Категория: Спорт                  │
├──────────────────────────────────────┤
│                                       │
│ 📝 Краткая ML-выжимка из 2-3          │
│    предложений, созданная             │
│    алгоритмом TextRank.               │
│                                       │
├──────────────────────────────────────┤
│ 🔗 [Перейти в источник →]            │
└──────────────────────────────────────┘
```

### Доступные категории

- 👗 Глянец
- 💊 Здоровье
- 🌍 Климат
- ⚔️ Конфликты
- 🎭 Культура
- 🔬 Наука
- 👥 Общество
- 🏛️ Политика
- ✈️ Путешествия
- ⚽ Спорт
- 💰 Экономика

---

## 🛠️ Установка (детальная)

### Требования

- **Go** 1.18+
- **Python** 3.8+
- **Git**
- **Windows PowerShell** 5.1+

### Шаг 1: Клонирование репозитория

```powershell
git clone <repository_url>
cd news_summarizer
```

### Шаг 2: Установка Go зависимостей

```powershell
cd backend
go mod download
go mod tidy
```

### Шаг 3: Установка Python зависимостей

**ML модуль:**
```powershell
cd ../ml
pip install -r requirements.txt
```

**Telegram парсер:**
```powershell
cd ../parser/tg_parser
python -m venv .venv
.\.venv\Scripts\activate
pip install -r requirements.txt
```

**Website парсер:**
```powershell
cd ../website_parser
python -m venv .venv
.\.venv\Scripts\activate
pip install -r requirements.txt
```

### Шаг 4: Настройка Telegram парсера

Создайте файл `parser/tg_parser/.env`:

```env
API_ID=your_telegram_api_id
API_HASH=your_telegram_api_hash
PHONE=+your_phone_number
```

Получить API credentials: https://my.telegram.org/apps

### Шаг 5: Первый запуск

```powershell
cd ../..
.\start_server.ps1
```

Откройте браузер: http://localhost:8081

---

## 💾 База данных

### Структура таблицы `news`

| Поле | Тип | Индекс | Описание |
|------|-----|--------|----------|
| `id` | INTEGER | PRIMARY KEY | Уникальный идентификатор |
| `title` | TEXT | - | Заголовок новости |
| `content` | TEXT | - | Полный текст статьи |
| `summary` | TEXT | - | **ML-выжимка (2-3 предложения)** |
| `source` | TEXT | INDEX | Источник (канал/сайт) |
| `url` | TEXT | **UNIQUE** | Ссылка на оригинал |
| `category` | TEXT | INDEX | Категория новости |
| `published_at` | DATETIME | INDEX | Дата публикации |
| `created_at` | DATETIME | - | Дата добавления в БД |

### Просмотр содержимого БД

```powershell
cd backend
go run inspector.go
```

**Пример вывода:**
```
Содержимое таблицы 'news':
ID: 1, Title: Заголовок, Source: BBC Russian, Category: Наука, Summary: Краткая выжимка...
ID: 2, Title: Другая новость, Source: Telegram, Category: Спорт, Summary: Еще выжимка...
```

### Настройки БД (оптимизация)

- **Journal Mode**: WAL (Write-Ahead Logging)
- **Busy Timeout**: 5000ms
- **Max Open Connections**: 1 (избегаем блокировок)

---

## 🤖 ML Суммаризация

### Алгоритм: TextRank

**Принцип работы:**
1. Разбивка текста на предложения
2. Построение графа схожести предложений
3. Ранжирование по алгоритму PageRank
4. Выбор топ-N предложений

**Особенности:**
- ✅ Экстрактивный метод (выбор из оригинального текста)
- ✅ Не требует обучения модели
- ✅ Быстрая работа
- ✅ Поддержка русского языка

### Настройка ML модуля

Файл `ml/summarize.py`:

```python
# Выбор алгоритма
from sumy.summarizers.text_rank import TextRankSummarizer
# Альтернативы: LexRankSummarizer, LsaSummarizer

# Язык
from sumy.nlp.tokenizers import Tokenizer
tokenizer = Tokenizer("russian")

# Длина выжимки (предложения)
sentence_count = 2  # или 3
```

---

## 🔄 Парсеры

### Telegram Parser (Telethon)

**Источники:**
- Настраиваются в `parser/tg_parser/mocks/channels_to_sub.json`

**Выходной файл:**
- `parser/tg_parser/mocks/export.json`

**Формат:**
```json
[
  {
    "title": "Заголовок сообщения",
    "content": "Полный текст",
    "url": "https://t.me/channel/12345",
    "channel": "Название канала"
  }
]
```

### Website Parser (BeautifulSoup)

**Источники:**
- BBC Russian (настроено по умолчанию)
- Добавьте свои в `parser/website_parser/script.py`

**Выходной файл:**
- `parser/website_parser/data/all_information.json`

**Формат:**
```json
[
  {
    "title": "Заголовок статьи",
    "content": "Полный текст",
    "url": "https://www.bbc.com/russian/articles/...",
    "category": "Наука"
  }
]
```

---

## 🐛 Решение проблем

### ❌ Сервер не запускается

**Проблема:** Порт 8081 занят

**Решение:**
```powershell
# Найти процесс
Get-NetTCPConnection -LocalPort 8081

# Убить процесс
Get-NetTCPConnection -LocalPort 8081 | ForEach-Object { 
    Stop-Process -Id $_.OwningProcess -Force 
}
```

### ❌ Парсеры не работают

**Telegram:**
1. Проверьте `.env` файл с credentials
2. Пройдите авторизацию (первый запуск)
3. Убедитесь, что установлены зависимости

**Website:**
1. Проверьте подключение к интернету
2. Убедитесь, что BeautifulSoup установлен
3. Проверьте логи в консоли сервера

### ❌ ML не генерирует выжимки

**Решение:**
```powershell
# Установите sumy
pip install sumy

# Проверьте Python
python --version

# Установите переменную окружения
$env:ML_PYTHON = "C:\path\to\python.exe"
```

### ❌ Кодировка в логах

**Проблема:** Иероглифы в PowerShell консоли

**Решение:** Это особенность Windows PowerShell с UTF-8.  
В браузере всё отображается корректно! ✅

---

## 📈 Производительность

| Операция | Время | Примечание |
|----------|-------|------------|
| Загрузка 20 новостей | < 100ms | Из БД |
| Парсинг + ML (полный цикл) | 1-2 мин | Фоново |
| Генерация 1 выжимки | ~2 сек | TextRank |
| API refresh (ответ) | < 50ms | 202 Accepted |

**Оптимизации:**
- ✅ SQLite WAL mode (параллельные чтения)
- ✅ Асинхронные задачи (goroutines)
- ✅ Timeout на ML операции (15 сек)
- ✅ Индексы на url, category, source

---

## 🤝 Разработка

### Добавление нового парсера

1. **Создайте скрипт** в `parser/your_parser/script.py`

2. **Формат вывода** (JSON):
```json
[
  {
    "title": "Заголовок",
    "content": "Полный текст",
    "url": "https://...",
    "category": "Категория"
  }
]
```

3. **Добавьте в Go** (`backend/main.go`):
```go
parsers := []string{
    "parser/tg_parser/parser_prototype.py",
    "parser/website_parser/script.py",
    "parser/your_parser/script.py", // ← Ваш парсер
}
```

### Запуск тестов

```powershell
# Backend тесты
cd backend
go test ./...

# Проверка БД
go run inspector.go
```

---

## 📝 Лицензия

MIT License - свободно используйте в своих проектах

---

## 👤 Автор

**Проект создан с помощью GitHub Copilot**  
Дата: 5 октября 2025 г.

---

## 🎯 Roadmap

### v1.1 (планируется)
- [ ] Docker контейнеризация
- [ ] Поддержка Docker Compose для быстрого развертывания
- [ ] Экспорт новостей в PDF
- [ ] Email уведомления

### v2.0 (в планах)
- [ ] Дополнительные источники (Reddit, Twitter/X)
- [ ] ML категоризация новостей (автоматическая)
- [ ] Персональные фильтры пользователя
- [ ] Sentiment analysis (тональность новостей)
- [ ] Графики и статистика

---

## ⭐ Благодарности

- **Sumy** - отличная библиотека для суммаризации
- **Telethon** - удобный клиент Telegram API
- **Gorilla Mux** - мощный роутер для Go
- **modernc.org/sqlite** - чистый Go драйвер SQLite

---

## 🚀 Начните сейчас!

```powershell
.\start_server.ps1
```

Затем откройте: **http://localhost:8081**

**Совет:** Добавьте в закладки для быстрого доступа! 🔖

---

**📧 Вопросы?** См. [USAGE_GUIDE.md](USAGE_GUIDE.md) или создайте Issue.

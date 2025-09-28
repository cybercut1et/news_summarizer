#!/bin/bash

# Скрипт для тестирования Backend API

BASE_URL="http://localhost:8080/api"

echo "=== Тестирование Backend API ==="

# 1. Создание пользователя
echo "1. Создание пользователя..."
USER_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "uuid": "test-user-123",
    "sources": ["lenta.ru", "rbc.ru", "tass.ru"]
  }')
echo "Ответ: $USER_RESPONSE"

# 2. Получение источников пользователя
echo -e "\n2. Получение источников пользователя..."
SOURCES_RESPONSE=$(curl -s -X GET $BASE_URL/users/test-user-123/sources)
echo "Ответ: $SOURCES_RESPONSE"

# 3. Добавление новости
echo -e "\n3. Добавление новости..."
NEWS_RESPONSE=$(curl -s -X POST $BASE_URL/news \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Тестовая новость",
    "content": "Это длинный текст новости для тестирования системы суммаризации. В данной новости рассказывается о важных событиях, которые произошли сегодня. Новость содержит множество деталей и подробностей, которые могут быть интересны читателям.",
    "source": "test-source",
    "url": "https://example.com/news/1",
    "published_at": "2024-01-01T12:00:00Z",
    "category": "технологии"
  }')
echo "Ответ: $NEWS_RESPONSE"

# 4. Получение списка новостей
echo -e "\n4. Получение списка новостей..."
NEWS_LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/news?limit=5")
echo "Ответ: $NEWS_LIST_RESPONSE"

# 5. Тестирование суммаризации
echo -e "\n5. Тестирование суммаризации..."
SUMMARIZE_RESPONSE=$(curl -s -X POST $BASE_URL/summarize \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Искусственный интеллект продолжает развиваться быстрыми темпами. Новые технологии машинного обучения позволяют создавать более точные модели предсказания. Компании активно внедряют AI в свои бизнес-процессы. Это приводит к значительному повышению эффективности работы. Однако существуют и определенные риски, связанные с автоматизацией."
  }')
echo "Ответ: $SUMMARIZE_RESPONSE"

echo -e "\n=== Тестирование завершено ==="
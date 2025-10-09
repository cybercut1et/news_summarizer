# Dockerfile для ML модуля
FROM python:3.9-slim

WORKDIR /app/ml

# Копируем requirements.txt
COPY requirements.txt .

# Устанавливаем Python зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Скачиваем NLTK данные
RUN python -c "import nltk; nltk.download('punkt')"

# Копируем Python скрипты
COPY . .

CMD ["python", "scripts/main.py"]
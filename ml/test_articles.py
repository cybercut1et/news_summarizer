# Загрузка текста
with open('test.txt', 'r', encoding='utf-8') as file:
    text = file.read()

# Разбиение текста на статьи, если нужно
articles = text.split('article_separator')
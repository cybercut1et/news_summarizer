import sumy
from sumy.parsers.plaintext import PlaintextParser
from sumy.nlp.tokenizers import Tokenizer
from sumy.summarizers.lsa import LsaSummarizer
from sumy.summarizers.text_rank import TextRankSummarizer

# Загрузка текста
with open('test.txt', 'r', encoding='utf-8') as file:
    text = file.read()

# Разбиение текста на статьи, если нужно
articles = text.split('article_separator')

# Функция для суммаризации текста
def summarize(text, sentence_count=2):
    # Для русского языка используем правильный токенизатор
    parser = PlaintextParser.from_string(text, Tokenizer("russian"))
    
    # Выбираем нужный суммаризатор
    summarizer = LsaSummarizer()  # Можно также использовать TextRankSummarizer()
    
    # Получаем резюме
    summary = summarizer(parser.document, sentence_count)
    return ' '.join(str(sentence) for sentence in summary)

# Применение суммаризации для каждой статьи
for article in articles:
    print("-----SUMMARY-----")
    print(summarize(article, 2))
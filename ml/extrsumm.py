import sumy
from sumy.parsers.plaintext import PlaintextParser
from sumy.nlp.tokenizers import Tokenizer
from sumy.summarizers.lsa import LsaSummarizer
from sumy.summarizers.text_rank import TextRankSummarizer
from test_articles import articles
import re

# Функция для суммаризации текста
def summarize(text, sentence_count=2):
    # Для русского языка используем правильный токенизатор
    parser = PlaintextParser.from_string(text, Tokenizer("russian"))
    
    # Выбираем нужный суммаризатор
    summarizer = LsaSummarizer()  # Можно также использовать TextRankSummarizer()
    
    # Получаем резюме
    summary = summarizer(parser.document, sentence_count)
    return ' '.join(str(sentence) for sentence in summary)

def clean_text(text):
    text = text.lower()
    text = re.sub(r'[^\sa-zA-Z0-9@\[\]]',' ',text) # Удаляет пунктцацию
    text = re.sub(r'\w*\d+\w*', '', text) # Удаляет цифры
    text = re.sub('\s{2,}', " ", text) # Удаляет ненужные пробелы
    return text


from classifier import classifier

# Применение суммаризации для каждой статьи
for article in articles:
    print("-----SUMMARY-----")
    summary = summarize(article, 2)
    print(classifier(summary))
    print(summary)
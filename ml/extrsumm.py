import sumy
from sumy.parsers.plaintext import PlaintextParser
from sumy.nlp.tokenizers import Tokenizer
from sumy.summarizers.lsa import LsaSummarizer
from sumy.summarizers.text_rank import TextRankSummarizer
import re
import json


from fasttext_classifier import classified_data



# Функция для суммаризации текста
def summarize(text, sentence_count):
    # Для русского языка используем правильный токенизатор
    parser = PlaintextParser.from_string(text, Tokenizer("russian"))
    
    # Выбираем нужный суммаризатор
    summarizer = TextRankSummarizer()  # Можно также использовать TextRankSummarizer()
    
    # Получаем резюме
    summary = summarizer(parser.document, sentence_count)
    return ' '.join(str(sentence) for sentence in summary)

if __name__ != "__extrsumm__":
    filtered_data = []
    # Применение суммаризации для каждой статьи
    for article in classified_data:
        if article["confidence"] > 0.9:
            article["content"] = summarize(article["content"], 1)
            filtered_data.append(article)
    with open("classified_test.json", "w", encoding="utf-8") as f:
        json.dump(filtered_data, f, ensure_ascii=False, indent=4)
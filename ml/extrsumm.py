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
    for channel in classified_data:
        channel_dict = {
            "channel_name": channel["channel_name"],
            "messages": []
        }
        for post in channel["messages"]:
            if post["confidence"] > 0.75:
                summarized_text = summarize(post["text"], 2)
                # Оставляем только нужные поля
                filtered_post = {
                    "text": summarized_text,
                    "date": post["date"],
                    "link": post["link"],
                    "category": post.get("category"),
                    "confidence": post.get("confidence")
                }
                channel_dict["messages"].append(filtered_post)
        filtered_data.append(channel_dict)
    with open("classified_test.json", "w", encoding="utf-8") as f:
        json.dump(filtered_data, f, ensure_ascii=False, indent=4)
import sumy
from sumy.parsers.plaintext import PlaintextParser
from sumy.nlp.tokenizers import Tokenizer
from sumy.summarizers.lsa import LsaSummarizer
from sumy.summarizers.text_rank import TextRankSummarizer
import re
import json
from pathlib import Path



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
def do_shit():
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
    output_path = Path(__file__).parent.parent / "filtered_data.json"
    with open(output_path, "w", encoding="utf-8") as f:
        json.dump(filtered_data, f, ensure_ascii=False, indent=4)

do_shit()
if __name__ == "__main__":
    import sys
    try:
        # Если скрипт вызывается с аргументом или через stdin
        if not sys.stdin.isatty():
            input_json = sys.stdin.read()
            try:
                data = json.loads(input_json)
                text = data.get("content", "")
                sentence_count = int(data.get("sentence_count", 2))
                summary_text = summarize(text, sentence_count)
                print(json.dumps({"summary": summary_text, "status": "ok"}, ensure_ascii=False))
            except Exception as e:
                print(json.dumps({"status": "error", "error": str(e)}))
        else:
            # Обычный запуск для генерации filtered_data.json
            do_shit()
    except Exception as e:
        print(json.dumps({"status": "error", "error": str(e)}))
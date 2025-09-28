import sys
import json
import sumy
from sumy.parsers.plaintext import PlaintextParser
from sumy.nlp.tokenizers import Tokenizer
from sumy.summarizers.lsa import LsaSummarizer
from sumy.summarizers.text_rank import TextRankSummarizer

def summarize_text(text, sentence_count=2):
    """
    Функция для суммаризации текста
    
    Args:
        text (str): Исходный текст для суммаризации
        sentence_count (int): Количество предложений в резюме
    
    Returns:
        str: Суммаризованный текст
    """
    try:
        # Для русского языка используем правильный токенизатор
        parser = PlaintextParser.from_string(text, Tokenizer("russian"))
        
        # Используем TextRank суммаризатор
        summarizer = TextRankSummarizer()
        
        # Получаем резюме
        summary = summarizer(parser.document, sentence_count)
        return ' '.join(str(sentence) for sentence in summary)
    except Exception as e:
        # В случае ошибки возвращаем первые 200 символов
        return text[:200] + "..." if len(text) > 200 else text

def main():
    """
    Основная функция для вызова из Go бэкенда
    Читает JSON из stdin, возвращает результат в stdout
    """
    try:
        # Читаем входные данные из stdin
        input_data = json.loads(sys.stdin.read())
        text = input_data.get('content', '')
        sentence_count = input_data.get('sentence_count', 2)
        
        # Выполняем суммаризацию
        summary = summarize_text(text, sentence_count)
        
        # Возвращаем результат
        result = {
            'summary': summary,
            'status': 'success'
        }
        print(json.dumps(result, ensure_ascii=False))
        
    except Exception as e:
        error_result = {
            'error': str(e),
            'status': 'error'
        }
        print(json.dumps(error_result, ensure_ascii=False))
        sys.exit(1)

if __name__ == "__main__":
    main()
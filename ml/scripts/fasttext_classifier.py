from huggingface_hub import hf_hub_download
import fasttext
import json
from pathlib import Path

class FastTextClassifierPipeline:
    def __init__(self, model_path):
        self.model = fasttext.load_model(model_path)

    def __call__(self, texts):
        if isinstance(texts, str):
            texts = [texts]

        results = []
        for text in texts:
            prediction = self.model.predict(text)
            label = prediction[0][0].replace("__label__", "")
            score = float(prediction[1][0])
            results.append({"label": label, "score": score})

        return results


def pipeline(task="text-classification", model=None):
    # Загрузка файла model.bin
    repo_id = "data-silence/fasttext-rus-news-classifier"
    model_file = hf_hub_download(repo_id=repo_id, filename="fasttext_news_classifier.bin")
    return FastTextClassifierPipeline(model_file)

# Создание классификатора
classifier = pipeline("text-classification")

def classify(tgfile, webfile):
    data = []
    with open(tgfile, 'r', encoding='utf-8') as file:
        tg_data = json.load(file)
    for channel in tg_data:
        for post in channel["messages"]:
            content = post["text"].replace('\n', ' ')
            classification = classifier(content)
            post["category"] = classification[0]["label"]
            post["confidence"] = classification[0]["score"]
        data.append(channel)
    with open(webfile, 'r', encoding='utf-8') as file:
        web_data = json.load(file)
    for site in web_data:
        for article in site["messages"]:
            content = article["text"].replace('\n', ' ')
            classification = classifier(content)
            article["category"] = classification[0]["label"]
            article["confidence"] = classification[0]["score"]
        data.append(site)
#        data.append({"site_name" : site_name, "messages" : [{"category": classification[0]["label"], "confidence": classification[0]["score"]}]})
    return data

tg_path = Path(__file__).parent.parent.parent / 'parser' / 'tg_parser' / 'mocks' / 'export.json'
web_path = Path(__file__).parent.parent.parent / 'parser' / 'website_parser' / 'data' / 'all_information.json'

classified_data = classify(str(tg_path), str(web_path))
# output_path = Path(__file__).parent / "classified_data.json"
# with open(output_path, "w", encoding="utf-8") as f:
#     json.dump(classified_data, f, ensure_ascii=False, indent=4)
# print(classified_data)
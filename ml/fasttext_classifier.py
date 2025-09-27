from huggingface_hub import hf_hub_download
import fasttext
import json

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

def classify(json_file_name):
    with open(json_file_name, 'r', encoding='utf-8') as file:
        data = json.load(file)
    for article in data:
        content = article["content"].replace('\n', ' ')
        classication = classifier(content)
        article["category"] = classifier(content)[0]["label"]
        article["confidence"] = classifier(content)[0]["score"]
    return data

classified_data = classify('test.json')
# print(classified_data)
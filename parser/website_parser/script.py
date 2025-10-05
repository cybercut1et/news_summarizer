import requests
from bs4 import BeautifulSoup
import json
import time

url = "https://www.bbc.com/russian"

headers = {
    "Accept": (
        "text/html,application/xhtml+xml,application/xml;q=0.9,"
        "image/avif,image/webp,image/apng,*/*;q=0.8,"
        "application/signed-exchange;v=b3;q=0.7"
    ),
    "User-Agent": (
        "Mozilla/5.0 (X11; Linux x86_64) "
        "AppleWebKit/537.36 (KHTML, like Gecko) "
        "Chrome/117.0.0.0 Safari/537.36"
    )
}

req = requests.get(url, headers=headers)
src1 = req.text

with open("index.html", "w") as file:
    file.write(src1)
with open("index.html") as file:
    src = file.read()

soup = BeautifulSoup(src, "lxml")
all_links = soup.find_all(class_="bbc-1i4ie53 e1d658bg0")

all_news_headers = {}
for item in all_links:
    item_text = item.text
    item_link = item.get("href")

    all_news_headers[item_text] = item_link

with open("all_news_headers.json", "w") as file:
                # Сохраняем все заголовки в json файл
    json.dump(all_news_headers, file, indent=4, ensure_ascii=False)

with open("all_news_headers.json") as file:
    all_news = json.load(file)

all_information = []
result = {
    "channel_name": url,
    "messages": []
}
count = 0

for category_name, category_href in all_news.items():
    if count >= 50:
        break
    # убираем пробелы и спецсимволы из имени файла
    safe_name = category_name.replace(" ", "_").replace("/", "_")

    # делаем ссылку абсолютной
    if category_href.startswith("/"):
        category_href = "https://www.bbc.com" + category_href

    print(f"[{count}] Парсим: {category_name}")

    try:
        # получаем HTML статьи
        req = requests.get(url=category_href, headers=headers)
        src2 = req.text

        # парсим страницу
        soup = BeautifulSoup(src2, "lxml")
        date_tag = soup.find("time")
        if date_tag and date_tag.get("datetime"):
            article_date = date_tag["datetime"]
        elif date_tag and date_tag.text.strip():
            article_date = date_tag.text.strip()
        else:
            article_date = None
        # ищем все <p> — обычно там основной текст
        paragraphs = soup.find_all("p")
        article_text = " ".join([p.get_text(strip=True) for p in paragraphs])

        # сохраняем в общий словарь
        result["messages"].append({
            "text": article_text,
            "date": article_date,
            "link": category_href
        })


        # небольшая задержка, чтобы не спамить сервер
        time.sleep(1)

        count += 1

    except Exception as e:
        print(f"Ошибка при обработке {category_href}: {e}")

all_information.append(result)
# сохраняем итоговый JSON
with open("data/all_information.json", "w", encoding="utf-8") as file:
    json.dump(all_information, file, indent=4, ensure_ascii=False)
import requests
from bs4 import BeautifulSoup
import json

# url = "https://www.bbc.com/russian"
#
# headers = {
#     "Accept": (
#         "text/html,application/xhtml+xml,application/xml;q=0.9,"
#         "image/avif,image/webp,image/apng,*/*;q=0.8,"
#         "application/signed-exchange;v=b3;q=0.7"
#     ),
#     "User-Agent": (
#         "Mozilla/5.0 (X11; Linux x86_64) "
#         "AppleWebKit/537.36 (KHTML, like Gecko) "
#         "Chrome/117.0.0.0 Safari/537.36"
#     )
# }
#
# req = requests.get(url, headers=headers)
# src = req.text
#print(src)

# with open("index.html") as file:
#     src = file.read()
#
# soup = BeautifulSoup(src, "lxml")
# all_links = soup.find_all(class_="bbc-1i4ie53 e1d658bg0")
#
# all_news_headers = {}
# for item in all_links:
#     item_text = item.text
#     item_link = item.get("href")
#
#     all_news_headers[item_text] = item_link
#
# with open("all_news_headers.json", "w") as file:
#     json.dump(all_news_headers, file, indent=4, ensure_ascii=False)

with open("all_news_headers.json") as file:
    all_news = json.load(file)

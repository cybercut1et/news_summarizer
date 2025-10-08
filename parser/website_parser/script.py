# fast_tass_parser_api.py
import asyncio
import aiohttp
import json
import os
from bs4 import BeautifulSoup
from aiohttp import ClientTimeout
from urllib.parse import urljoin

BASE_URL = "https://tass.ru"
API_URL = "https://tass.ru/rss/v2.xml"  # официальный RSS-канал
HEADERS = {
    "User-Agent": (
        "Mozilla/5.0 (X11; Linux x86_64) "
        "AppleWebKit/537.36 (KHTML, like Gecko) "
        "Chrome/117.0.0.0 Safari/537.36"
    )
}

MAX_ARTICLES = 10
CONCURRENCY = 5
TIMEOUT = ClientTimeout(total=25)

async def fetch_text(session: aiohttp.ClientSession, url: str, retries: int = 3) -> str:
    for attempt in range(retries):
        try:
            async with session.get(url, headers=HEADERS, timeout=TIMEOUT) as resp:
                resp.raise_for_status()
                return await resp.text()
        except (aiohttp.ClientError, asyncio.TimeoutError) as e:
            if attempt < retries - 1:
                print(f"⚠️ Ошибка {type(e).__name__} при загрузке {url}, попытка {attempt + 1}/{retries}...")
                await asyncio.sleep(2 * (attempt + 1))
                continue
            else:
                raise


async def parse_rss(session: aiohttp.ClientSession):
    """Парсим RSS TASS и получаем заголовки и ссылки"""
    xml = await fetch_text(session, API_URL)
    soup = BeautifulSoup(xml, "xml")
    items = soup.find_all("item")

    articles = []
    for it in items[:MAX_ARTICLES]:
        title = it.title.get_text(strip=True)
        link = it.link.get_text(strip=True)
        articles.append((title, link))
    return articles

def extract_article(soup: BeautifulSoup):
    """Извлекаем дату и основной текст статьи Tass.ru"""
    # ==== 1. Попробуем найти тег <time> ====
    article_date = None
    time_tag = soup.find("time")
    if time_tag:
        if time_tag.get("datetime"):
            article_date = time_tag["datetime"]
        elif time_tag.get_text(strip=True):
            article_date = time_tag.get_text(strip=True)

    # ==== 2. Попробуем meta-теги ====
    if not article_date:
        for key in [
            ("meta", {"property": "article:published_time"}),
            ("meta", {"name": "pubdate"}),
            ("meta", {"property": "og:pubdate"}),
            ("meta", {"property": "og:updated_time"}),
            ("meta", {"name": "Last-Modified"}),
        ]:
            tag = soup.find(*key)
            if tag and tag.get("content"):
                article_date = tag["content"]
                break

    # ==== 3. Попробуем дата в тексте (если ничего не нашли) ====
    if not article_date:
        # пример формата: "8 октября 2025, 14:12"
        possible_time = soup.select_one("span.article__header__date")
        if possible_time:
            article_date = possible_time.get_text(strip=True)

    # ==== 4. Извлекаем основной текст ====
    text_blocks = soup.select("div.article__text p")
    if not text_blocks:
        text_blocks = soup.select("article p")

    article_text = " ".join(
        p.get_text(" ", strip=True)
        for p in text_blocks
        if p.get_text(strip=True)
    )

    return {"date": article_date, "text": article_text}



async def fetch_article(session, title, url, sem):
    async with sem:
        try:
            html = await fetch_text(session, url)
            soup = BeautifulSoup(html, "lxml")
            data = extract_article(soup)
            return {
                "header": title,
                "text": data["text"],
                "date": data["date"],
                "link": url
            }
        except Exception as e:
            return {
                "header": title,
                "text": "",
                "date": None,
                "link": url,
                "error": str(e)
            }


async def main():
    os.makedirs("data", exist_ok=True)
    connector = aiohttp.TCPConnector(limit_per_host=CONCURRENCY, ssl=False)
    sem = asyncio.Semaphore(CONCURRENCY)

    async with aiohttp.ClientSession(connector=connector) as session:
        articles = await parse_rss(session)
        print(f"🔍 Найдено {len(articles)} статей в RSS. Загружаем тексты...")

        tasks = [fetch_article(session, t, l, sem) for t, l in articles]
        results = await asyncio.gather(*tasks)

    all_information = [{
        "channel_name": BASE_URL,
        "messages": results
    }]

    out_path = "data/all_information.json"
    with open(out_path, "w", encoding="utf-8") as f:
        json.dump(all_information, f, ensure_ascii=False, indent=4)

    print(f"✅ Готово! Сохранено {len(results)} статей в {out_path}")


if __name__ == "__main__":
    asyncio.run(main())

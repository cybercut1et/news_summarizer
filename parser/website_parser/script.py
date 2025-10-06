# fast_bbc_parser.py
import asyncio
import aiohttp
from aiohttp import ClientTimeout
from bs4 import BeautifulSoup
import json
import os
from typing import List, Dict, Tuple

BASE_URL = "https://www.bbc.com"
CHANNEL_URL = f"{BASE_URL}/russian"

HEADERS = {
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

# Сколько статей тянуть максимум
MAX_ARTICLES = 20
# Сколько запросов одновременно
CONCURRENCY = 8
# Таймауты на все операции
TIMEOUT = ClientTimeout(total=12)


def absolutize(href: str) -> str:
    if not href:
        return href
    return href if href.startswith("http") else f"{BASE_URL}{href}"


async def fetch_text(session: aiohttp.ClientSession, url: str, *, retries: int = 2) -> str:
    """
    Быстрый и надёжный GET с 2 повторами на сетевые ошибки/5xx.
    """
    for attempt in range(retries + 1):
        try:
            async with session.get(url, headers=HEADERS, timeout=TIMEOUT) as resp:
                # BBC корректно шлёт gzip/deflate; aiohttp это распакует сам.
                if resp.status >= 400:
                    raise aiohttp.ClientResponseError(
                        request_info=resp.request_info, history=resp.history,
                        status=resp.status, message=f"HTTP {resp.status}"
                    )
                return await resp.text()
        except (aiohttp.ClientError, asyncio.TimeoutError) as e:
            if attempt == retries:
                raise
            # Небольшая экспоненциальная пауза между повторами
            await asyncio.sleep(0.4 * (attempt + 1))


async def parse_index(session: aiohttp.ClientSession) -> List[Tuple[str, str]]:
    """
    Парсим главную, достаём пары (заголовок, ссылка).
    Используем твой селектор класса, плюс небольшой запасной вариант.
    """
    html = await fetch_text(session, CHANNEL_URL)
    soup = BeautifulSoup(html, "lxml")

    # 1) Твой селектор
    items = soup.select(".bbc-1i4ie53.e1d658bg0")
    # 2) Запасной путь: любые ссылки внутри заголовочных блоков
    if not items:
        items = soup.select("a[href]")

    pairs: List[Tuple[str, str]] = []
    seen_links = set()

    for a in items:
        text = a.get_text(strip=True)
        href = a.get("href")
        if not text or not href:
            continue

        url = absolutize(href)
        # Немного фильтрации, чтобы не брать сервисные страницы
        if "/russian" not in url:
            continue

        # Простейшая дедупликация по ссылке
        if url in seen_links:
            continue
        seen_links.add(url)

        pairs.append((text, url))
        if len(pairs) >= MAX_ARTICLES * 2:
            # берём с запасом, потом обрежем
            break

    # Сохраним порядок и ужмём до MAX_ARTICLES
    return pairs[:MAX_ARTICLES]


def extract_article(soup: BeautifulSoup) -> Dict:
    """
    Быстрый и «достаточно хороший» извлекатель:
    - дата из <time datetime|text>
    - текст: <main> -> все <p>, иначе все <p>
    """
    # Дата/время
    article_date = None
    time_tag = soup.find("time")
    if time_tag and time_tag.get("datetime"):
        article_date = time_tag["datetime"]
    elif time_tag and time_tag.get_text(strip=True):
        article_date = time_tag.get_text(strip=True)

    # Текст
    main = soup.find("main")
    paragraphs = (main or soup).find_all("p")
    article_text = " ".join(p.get_text(strip=True) for p in paragraphs if p.get_text(strip=True))

    return {"date": article_date, "text": article_text}


async def fetch_and_parse_article(
    session: aiohttp.ClientSession, title: str, url: str, sem: asyncio.Semaphore
) -> Dict:
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
            # Возвращаем заглушку вместо падения всего парсинга
            return {
                "header": title,
                "text": "",
                "date": None,
                "link": url,
                "error": f"{type(e).__name__}: {e}"
            }


async def main():
    os.makedirs("data", exist_ok=True)

    connector = aiohttp.TCPConnector(limit_per_host=CONCURRENCY, ssl=False)  # ssl=False ускоряет хэндшейк на некоторых системах
    sem = asyncio.Semaphore(CONCURRENCY)

    async with aiohttp.ClientSession(connector=connector) as session:
        pairs = await parse_index(session)

        tasks = [
            asyncio.create_task(fetch_and_parse_article(session, title, link, sem))
            for title, link in pairs
        ]
        messages = await asyncio.gather(*tasks)

    result = {
        "channel_name": CHANNEL_URL,
        "messages": messages
    }

    all_information = [result]

    out_path = "data/all_information.json"
    with open(out_path, "w", encoding="utf-8") as f:
        json.dump(all_information, f, indent=4, ensure_ascii=False)

    print(f"Готово! Сохранено: {out_path} — статей: {len(messages)}")


if __name__ == "__main__":
    asyncio.run(main())

import sqlite3
import csv
from datetime import datetime

DB = './news.db'
OUT = './news_dump_100.csv'

conn = sqlite3.connect(DB)
conn.row_factory = sqlite3.Row
cur = conn.cursor()

query = '''SELECT id, title, content, summary, source, url, published_at, created_at, category
FROM news
ORDER BY published_at DESC
LIMIT 100'''

rows = cur.execute(query).fetchall()

with open(OUT, 'w', encoding='utf-8-sig', newline='') as f:
    writer = csv.writer(f)
    writer.writerow(['id','title','content','summary','source','url','published_at','created_at','category'])
    for r in rows:
        # Convert datetime-like fields to ISO strings if possible
        published = r['published_at']
        created = r['created_at']
        try:
            # sqlite may return strings; keep as-is otherwise
            if isinstance(published, str):
                published_val = published
            else:
                published_val = str(published)
        except Exception:
            published_val = ''
        try:
            if isinstance(created, str):
                created_val = created
            else:
                created_val = str(created)
        except Exception:
            created_val = ''
        writer.writerow([
            r['id'],
            r['title'] if r['title'] is not None else '',
            r['content'] if r['content'] is not None else '',
            r['summary'] if r['summary'] is not None else '',
            r['source'] if r['source'] is not None else '',
            r['url'] if r['url'] is not None else '',
            published_val,
            created_val,
            r['category'] if r['category'] is not None else ''
        ])

conn.close()
print(f"Exported {len(rows)} rows to {OUT}")

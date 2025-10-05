import sqlite3
import csv

DB = './news.db'
OUT = './news_dump_100_fixed.csv'

def decode_field(b):
    if b is None:
        return ''
    if isinstance(b, str):
        return b
    # b is bytes
    try:
        return b.decode('utf-8')
    except Exception:
        try:
            return b.decode('cp1251')
        except Exception:
            return b.decode('utf-8', errors='replace')

conn = sqlite3.connect(DB)
# return raw bytes for text columns
conn.text_factory = bytes
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
        # decode fields
        idv = r[0]
        title = decode_field(r[1])
        content = decode_field(r[2])
        summary = decode_field(r[3])
        source = decode_field(r[4])
        url = decode_field(r[5])
        published = r[6] if r[6] is not None else ''
        created = r[7] if r[7] is not None else ''
        category = decode_field(r[8])
        writer.writerow([idv, title, content, summary, source, url, published, created, category])

conn.close()
print(f"Exported {len(rows)} rows to {OUT}")

import sqlite3
import textwrap

DB = './news.db'

def trunc(s, n=80):
    if s is None:
        return ''
    s = str(s)
    return (s[:n-3] + '...') if len(s) > n else s

conn = sqlite3.connect(DB)
conn.row_factory = sqlite3.Row
cur = conn.cursor()

print('\nTABLE: news (first 20 rows)')
print('-' * 120)
print(f"{'id':>3}  {'title':30}  {'source':15}  {'url':45}  {'summary'}")
print('-' * 120)
for row in cur.execute('SELECT id, title, source, url, summary FROM news ORDER BY published_at DESC LIMIT 20'):
    print(f"{row['id']:>3}  {trunc(row['title'],30):30}  {trunc(row['source'],15):15}  {trunc(row['url'],45):45}  {trunc(row['summary'],60)}")

print('\nTABLE: users (first 20 rows)')
print('-' * 80)
print(f"{'id':>3}  {'uuid':36}  {'sources'}")
print('-' * 80)
for row in cur.execute('SELECT id, uuid, sources FROM users LIMIT 20'):
    print(f"{row['id']:>3}  {trunc(row['uuid'],36):36}  {trunc(row['sources'],40)}")

conn.close()

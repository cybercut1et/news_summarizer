from telethon import TelegramClient
import os
from dotenv import load_dotenv
import asyncio
import time
import json
load_dotenv()

api_id = os.getenv('API_ID')
api_hash = os.getenv('API_HASH')

client = TelegramClient('tg_session', api_id, api_hash, system_version='4.16.30-vxhello', device_model='Tecno TECNO CAMON 20 PRO')

async def main():
    dialogs = await client.get_dialogs()
    for dialog in dialogs:
        if dialog.title == 'Взял Мяч':
            messages = await client.get_messages(dialog, limit=10)
            for message in messages:
                print(message.text)
                print('==' * 50)
                time.sleep(1)

with client:
    client.loop.run_until_complete(main())

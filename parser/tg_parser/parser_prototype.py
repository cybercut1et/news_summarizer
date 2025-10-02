from telethon import TelegramClient
import os
from dotenv import load_dotenv
import asyncio
from datetime import datetime
import json
import pytz

load_dotenv()
# получаем api приложения
api_id = os.getenv('API_ID')
api_hash = os.getenv('API_HASH')


time_period = json.load(open('time_period.json', 'r'))['hours']
# получаем текущее время в минутах и текущий день
current_date = datetime.now().date()
current_time = datetime.now().time()
minutes_now = current_time.hour * 60 + current_time.minute
current_day = current_date.day
# для вывода по часам нужно знать московское время
moscow_tz = pytz.timezone('Europe/Moscow')

# создаем юзербота
client = TelegramClient('tg_session', api_id, api_hash, system_version='4.16.30-vxhello', device_model='Tecno TECNO CAMON 20 PRO')

async def parse_channels():
    dialogs = await client.get_dialogs() # получаем множество диалогов
    export_data = [] # итоговый файл

    channels_to_parse = json.load(open('channels.json', 'r', encoding='utf-8'))['channels']

    # среди диалогов ищем нужный нам
    for dialog in dialogs:
        if dialog.title in channels_to_parse:
            chat_id = dialog.entity.username # юзер тгк
            messages = await client.get_messages(dialog, limit=20) # пока что достаем 10 последних новостей, чтобы не наглеть
            
            # в export_data будет название канала и последние новости
            parsed_data = {'channel_name': dialog.title, 'messages': []}
            for message in messages:
                # считаем время и дату публикации и разбиваем на части
                message_publication_date_time = message.date.astimezone(moscow_tz).isoformat().split('+')[0] # здесь сплит так как в iso есть часовой пояс, после +
                message_publication_time = message_publication_date_time.split('T')[1] # время публикации
                message_publication_date = message_publication_date_time.split('T')[0] # дата публикации

                # считаем день и минуты публикации
                minutes_publication = int(message_publication_time.split(':')[0]) * 60 + int(message_publication_time.split(':')[1])
                day_publication = int(message_publication_date.split('-')[2])

                # считаю разницу между временем на компьютере и выпуском поста, если она меньше нужной, то не вывожу
                request_publication_diff = (minutes_now - minutes_publication) if current_day == day_publication else minutes_now + (1440 - minutes_publication)
                
                if (request_publication_diff <= time_period * 60) and (message.text): # здесь идет проверка на message.text, чтобы не было пустых постов
                    # в message_data текст публикации, датаа выпуска публикации и ссылка на публикацию
                    message_data = { 
                        'text': message.text,
                        'date': message_publication_time.replace('T', ' '),
                        'link' : f'https://t.me/{chat_id}/{message.id}' if chat_id is not None else None
                    }
                    parsed_data['messages'].append(message_data)
            export_data.append(parsed_data)
    return export_data

async def main():
    async with client:
        export_data = await parse_channels()
        with open('export.json', 'w', encoding='utf-8') as f:
            json.dump(export_data, f, ensure_ascii=False, indent=4)

if __name__ == '__main__':
    asyncio.run(main())

from telethon import TelegramClient
import os
from dotenv import load_dotenv
import asyncio
from datetime import datetime
import json
import pytz
load_dotenv()

time_period = int(input('За сколько часов хотите получить новости: '))

# получаем api приложения
api_id = os.getenv('API_ID')
api_hash = os.getenv('API_HASH')

current_time = datetime.now()
# для вывода по часам нужно знать московское время
moscow_tz = pytz.timezone('Europe/Moscow')

# создаем юзербота
client = TelegramClient('tg_session', api_id, api_hash, system_version='4.16.30-vxhello', device_model='Tecno TECNO CAMON 20 PRO')

async def main():
    # получаем множество диалогов
    dialogs = await client.get_dialogs()



    # среди диалогов ищем нужный нам(для примера канал Взял Мяч)
    for dialog in dialogs:
        if dialog.title == 'Взял Мяч':
            chat_id = dialog.entity.username
            messages = await client.get_messages(dialog, limit=10) # пока что достаем 10 последних новостей, чтобы не наглеть
            
            # в export_data будет название канала и последние новости
            export_data = {'channel_name': dialog.title, 'messages': []}
            for message in messages:
                message_publication_time = message.date.astimezone(moscow_tz).isoformat()

                # считаю разницу между временем на компьютере и выпуском поста, если она меньше нужной, то не вывожу
                request_publication_diff = abs(int(message_publication_time.split('T')[1].split('+')[0].split(':')[0]) - int(current_time.time().hour))
                
                if request_publication_diff <= time_period:
                    # в message_data текст публикации, датаа выпуска публикации и ссылка на публикацию
                    message_data = { 
                        'text': message.text,
                        'date': message_publication_time,
                        'link' : f'https://t.me/{chat_id}/{message.id}'
                    }
                    export_data['messages'].append(message_data)

            filename = f'got_ball_messages.json'
            with open(filename, 'w', encoding='utf-8') as f:
                json.dump(export_data, f, ensure_ascii=False)
            

with client:
    client.loop.run_until_complete(main())

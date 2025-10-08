from telethon.tl.functions.channels import JoinChannelRequest 
from telethon import TelegramClient
import os
from dotenv import load_dotenv
import asyncio
from datetime import datetime
import json
import pytz
from pathlib import Path

CHANNELS_TO_SUB = Path('mocks/channels_to_sub.json')
CHANNELS = Path('mocks/channels.json')
EXPORT = Path('mocks/export.json')
TIME_PERIOD = Path('mocks/time_period.json')

load_dotenv()
# получаем api приложения
api_id = os.getenv('API_ID')
api_hash = os.getenv('API_HASH')


time_period = json.load(open(TIME_PERIOD, 'r'))['hours']
# получаем текущее время в минутах и текущий день
current_date = datetime.now().date()
current_time = datetime.now().time()
minutes_now = current_time.hour * 60 + current_time.minute
current_day = current_date.day
# для вывода по часам нужно знать московское время
moscow_tz = pytz.timezone('Europe/Moscow')

# создаем юзербота
client = TelegramClient('tg_session', api_id, api_hash, system_version='4.16.30-vxhello', device_model='Tecno TECNO CAMON 20 PRO')

async def parse_channels(subbed_channels=[]):
    dialogs = await client.get_dialogs() # получаем множество диалогов
    export_data = [] # итоговый файл

    # смотрим каналы которые нужно считать, и добавляем туда каналы на которые мы подписались
    channels_to_parse = json.load(open(CHANNELS, 'r', encoding='utf-8'))['channels']
    channels_to_parse.extend(subbed_channels)

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

# функция для подписки на каналы и добавления их в список для парсинга
async def sub_to_channel(channel_name):
    try:
        # Получаем информацию о канале и подписываемся
        channel_entity = await client.get_entity(channel_name)
        await client(JoinChannelRequest(channel_entity))
        
        # Получаем название канала
        channel_title = channel_entity.title
        
        # Добавляем канал в список для парсинга
        try:
            # Читаем текущий список каналов для парсинга
            with open(CHANNELS, 'r', encoding='utf-8') as f:
                channels_data = json.load(f)
            
            # Добавляем новый канал если его еще нет
            if channel_title not in channels_data['channels']:
                channels_data['channels'].append(channel_title)
                with open(CHANNELS, 'w', encoding='utf-8') as f:
                    json.dump(channels_data, f, ensure_ascii=False, indent=4)
                print(f"Успешно подписались на канал '{channel_title}' и добавили его в список для парсинга")
            else:
                print(f"Канал '{channel_title}' уже есть в списке для парсинга")
            
            # Удаляем канал из списка для подписки
            with open(CHANNELS_TO_SUB, 'r', encoding='utf-8') as f:
                channels_to_sub_data = json.load(f)
            
            if channel_name in channels_to_sub_data['channels']:
                channels_to_sub_data['channels'].remove(channel_name)
                with open(CHANNELS_TO_SUB, 'w', encoding='utf-8') as f:
                    json.dump(channels_to_sub_data, f, ensure_ascii=False, indent=4)
                print(f"Канал '{channel_name}' удален из списка для подписки")
                
        except Exception as e:
            print(f"Подписались на канал '{channel_title}', но произошла ошибка при обновлении списков: {str(e)}")
            
    except Exception as e:
        print(f"Ошибка при подписке на канал {channel_name}: {str(e)}")
        raise e

async def main():
    async with client:
        # считываем каналы из файла и подписываемся на них
        channels_to_sub = json.load(open(CHANNELS_TO_SUB, 'r', encoding='utf-8'))['channels']
        for channel in channels_to_sub:
            await sub_to_channel(channel)


        export_data = await parse_channels()

        with open(EXPORT, 'w', encoding='utf-8') as f:
            json.dump(export_data, f, ensure_ascii=False, indent=4)

if __name__ == '__main__':
    asyncio.run(main())

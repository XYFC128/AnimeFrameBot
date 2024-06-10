from telethon import TelegramClient
from telethon.sessions import StringSession


api_id = int(input("api_id: "))
api_hash = input("api_hash: ")


with TelegramClient(StringSession(), api_id, api_hash) as client:
    print("Session string:", client.session.save())

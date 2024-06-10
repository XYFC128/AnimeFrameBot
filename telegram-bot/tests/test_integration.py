import pytest
import configparser
from telethon import TelegramClient
from telethon.sessions import StringSession
from telethon.tl.custom.message import Message
import subprocess
import signal

from telethon.tl.custom.conversation import Conversation


@pytest.fixture(autouse=True, scope="session")
def bot():
    # """Start bot to be tested."""
    process = subprocess.Popen(['python', 'src/main.py', 'config.ini'])
    yield
    process.send_signal(signal.SIGINT)
    

@pytest.fixture(scope="session")
async def telegram_client():
    """Connect to Telegram user for testing."""
    config = configparser.ConfigParser()
    config.read("config.ini")
    api_id = int(config['Telegram API']['api_id'])
    api_hash = config['Telegram API']['api_hash']
    session_str = config['Telegram API']['session_string']

    client = TelegramClient(
        StringSession(session_str), api_id, api_hash, sequential_updates=True
    )
    await client.connect()
    await client.get_me()
    await client.get_dialogs()

    yield client

    await client.disconnect()
    await client.disconnected


@pytest.mark.asyncio
async def test_hello(telegram_client):
    async with telegram_client.conversation(
        '@AnimeFrameBot', timeout=10, max_messages=10000
    ) as conv:
        conv: Conversation

        await conv.send_message("/start")
        res: Message = await conv.get_response()  # Welcome message
        assert res.text == "I'm a bot, please talk to me!"


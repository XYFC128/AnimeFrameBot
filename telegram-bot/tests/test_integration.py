import pytest
import configparser
from telethon import TelegramClient
from telethon.sessions import StringSession
from telethon.tl.custom.message import Message
import subprocess
import signal
import configparser

from telethon.tl.custom.conversation import Conversation


@pytest.fixture(autouse=True, scope="session")
def bot():
    # """Start bot to be tested."""
    process = subprocess.Popen(['python', 'src/main.py', 'config.ini'])
    yield
    process.send_signal(signal.SIGINT)

@pytest.fixture(scope="session")
def config():
    config = configparser.ConfigParser()
    config.read('config.ini')
    return config


@pytest.fixture(scope="session")
async def telegram_client(config):
    """Connect to Telegram user for testing."""
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


@pytest.fixture(scope="session")
async def conv(telegram_client, config):
    async with telegram_client.conversation(
        config['Bot']['name'], timeout=10, max_messages=10000
    ) as conv:
        conv: Conversation

        await conv.send_message("/start")
        await conv.get_response()  # Welcome message
        yield conv
        await conv.mark_read()


@pytest.mark.asyncio
async def test_start(conv: Conversation):
    await conv.send_message("/start")
    res: Message = await conv.get_response()
    assert res.text == "I'm a bot, please talk to me!"


@pytest.mark.asyncio
async def test_help(conv: Conversation):
    from src.main import HELP_TEXT
    await conv.send_message("/help")
    res: Message = await conv.get_response()
    assert res.text == HELP_TEXT.strip()

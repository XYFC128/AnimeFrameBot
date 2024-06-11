from telethon import TelegramClient
from telethon.sessions import StringSession
from telethon.tl.custom.conversation import Conversation
from telethon.tl.custom.message import Message
from telethon.tl.types import MessageMediaPhoto
import configparser
import configparser
import pytest
import signal
import subprocess


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


@pytest.mark.asyncio
async def test_frame_empty(conv: Conversation):
    await conv.send_message("/frame")
    res: Message = await conv.get_response()
    assert res.text == "Please provide a text to frame like this: /frame {text}"


@pytest.mark.asyncio
async def test_frame_bad_frame_num(conv: Conversation):
    await conv.send_message("/frame hello -1")
    res: Message = await conv.get_response()
    assert res.text == "Please provide a valid frame number."

    await conv.send_message("/frame hello 0")
    res: Message = await conv.get_response()
    assert res.text == "Please provide a valid frame number."

    await conv.send_message("/frame hello 100")
    res: Message = await conv.get_response()
    assert res.text == "I can only provide at most 10 frames at once."


@pytest.mark.asyncio
async def test_frame_normal(conv: Conversation):
    await conv.send_message("/frame 爽")
    res: Message = await conv.get_response()
    assert isinstance(res.media, MessageMediaPhoto) or res.text == "Sorry, I'm unable to find any frame."


@pytest.mark.asyncio
async def test_reandom_bad_frame_num(conv: Conversation):
    await conv.send_message("/random -1")
    res: Message = await conv.get_response()
    assert res.text == "Please provide a valid frame number."

    await conv.send_message("/random 0")
    res: Message = await conv.get_response()
    assert res.text == "Please provide a valid frame number."

    await conv.send_message("/random 100")
    res: Message = await conv.get_response()
    assert res.text == "I can only provide at most 10 frames at once."


@pytest.mark.asyncio
async def test_random_normal(conv: Conversation):
    await conv.send_message("/random")
    res: Message = await conv.get_response()
    assert isinstance(res.media, MessageMediaPhoto) or res.text == "Sorry, I'm unable to find any frame."


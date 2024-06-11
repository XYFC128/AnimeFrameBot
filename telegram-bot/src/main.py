from telegram import InputMediaPhoto, Update
from telegram.ext import ApplicationBuilder, ContextTypes, CommandHandler, MessageHandler, filters
import asyncio
import configparser
import logging
import os
import requests
import urllib.parse

FRAME_NUMBER = 1
API_URL = "http://localhost:8763"
TMP_DIR = "/tmp/AnimeFrameBot"
BOT_NAME = None
HELP_TEXT = """
I am a bot that can get frames you want from an anime with the text you provide.
Here are the commands you can use:
/help - Shows help message
/start - Starts the bot
/frame {query} {N} - Gets N (in range [1, 10], default: 1) frames with similar subtitle as the query
/random {N} - Gets random N (in range [1, 10], default: 1) frames

You can also send me an image and I will upload it to the server,
but be sure to provide a caption for the image if you upload it with compression.
"""

async def help_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(chat_id=update.effective_chat.id, text=HELP_TEXT)


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(chat_id=update.effective_chat.id, text="I'm a bot, please talk to me!")


async def frame(update: Update, context: ContextTypes.DEFAULT_TYPE):
    if not context.args:
        await context.bot.send_message(chat_id=update.effective_chat.id, text="Please provide a text to frame like this: /frame {text}")
        return
    
    text = context.args[0]
    if len(context.args) > 1:
        frame_number = context.args[1]
        if not frame_number.isdigit() or int(frame_number) == 0:
            await context.bot.send_message(chat_id=update.effective_chat.id, text="Please provide a valid frame number.")
            return
        elif int(frame_number) > 10:
            await context.bot.send_message(chat_id=update.effective_chat.id, text="I can only provide at most 10 frames at once.")
            return
    else:
        frame_number = FRAME_NUMBER
    text = urllib.parse.quote(text)
    url = f"{API_URL}/frame/fuzzy/{text}/{frame_number}"
    
    try:
        response = requests.get(url)
        frames = response.json()

        media_group = []
        for frame in frames:
            image_url = f"{API_URL}/frame/{urllib.parse.quote(frame['name'])}"
            image = requests.get(image_url).content
            media = InputMediaPhoto(image, caption=frame['subtitle'])
            media_group.append(media)
        if len(media_group) > 0:
            await update.effective_chat.send_media_group(media_group)
            return

    except Exception as e:
        logging.error(f"/frame failed: {e}")

    await context.bot.send_message(chat_id=update.effective_chat.id, text="Sorry, I'm unable to find any frame.")


async def random(update: Update, context: ContextTypes.DEFAULT_TYPE):
    if len(context.args) > 0:
        frame_number = context.args[0]
        if not frame_number.isdigit() or int(frame_number) == 0:
            await context.bot.send_message(chat_id=update.effective_chat.id, text="Please provide a valid frame number.")
            return
        elif int(frame_number) > 10:
            await context.bot.send_message(chat_id=update.effective_chat.id, text="I can only provide at most 10 frames at once.")
            return
    else:
        frame_number = FRAME_NUMBER

    url = f"{API_URL}/frame/random/{frame_number}"

    try:
        response = requests.get(url)
        frames = response.json()

        media_group = []
        for frame in frames:
            image_url = f"{API_URL}/frame/{urllib.parse.quote(frame['name'])}"
            image = requests.get(image_url).content
            media = InputMediaPhoto(image, caption=frame['subtitle'])
            media_group.append(media)
        if len(media_group) > 0:
            await update.effective_chat.send_media_group(media_group)
            return

    except Exception as e:
        logging.error(f"/random failed: {e}")

    await context.bot.send_message(chat_id=update.effective_chat.id, text="Sorry, I'm unable to find any frame.")


async def handle_smart_reply(update: Update, context: ContextTypes.DEFAULT_TYPE):
    logging.info(f"Smart reply for: {update.message.text}")
    text = update.message.text
    text = urllib.parse.quote(text)
    url = f"{API_URL}/frame/exact/{text}/{1}"
    
    try:
        response = requests.get(url)
        frames = response.json()

        media_group = []
        for frame in frames:
            image_url = f"{API_URL}/frame/{urllib.parse.quote(frame['name'])}"
            image = requests.get(image_url).content
            media = InputMediaPhoto(image, caption=frame['subtitle'])
            media_group.append(media)
        if len(media_group) > 0:
            await update.effective_chat.send_media_group(media_group, reply_to_message_id=update.message.id)

    except Exception as e:
        pass


async def upload(update: Update, context: ContextTypes.DEFAULT_TYPE, file_path: str):
    url = f"{API_URL}/frame"
    files = {'image': open(file_path, 'rb')}
    try:
        response = requests.post(url, files=files)
        if response.status_code == 201:
            await context.bot.send_message(chat_id=update.effective_chat.id, text=f"Image {file_path.split('/')[-1]} uploaded successfully")
        else:
            await context.bot.send_message(chat_id=update.effective_chat.id, text=f"Failed to upload image: {response.text}")
    except requests.RequestException as e:
        await context.bot.send_message(chat_id=update.effective_chat.id, text=f"upload failed: {e}")
    finally:
        files['image'].close()
        os.remove(file_path)


async def image_file_downloader(update: Update, context: ContextTypes.DEFAULT_TYPE):
    caption = update.message.caption
    if not caption:
        file_name = update.message.document.file_name
    else:
        file_name = caption + '.jpg'
    file_path = os.path.join(TMP_DIR, file_name)
    
    if not os.path.exists(TMP_DIR):
        os.makedirs(TMP_DIR)
    file = await context.bot.get_file(update.message.document.file_id)
    await file.download_to_drive(file_path)

    await upload(update, context, file_path)


async def image_downloader(update: Update, context: ContextTypes.DEFAULT_TYPE):
    caption = update.message.caption
    if not caption:
        await context.bot.send_message(chat_id=update.effective_chat.id, text="Please provide a caption for the image")
    file_path = os.path.join(TMP_DIR, caption + '.jpg')

    if not os.path.exists(TMP_DIR):
        os.makedirs(TMP_DIR)
    file = await context.bot.get_file(update.message.photo[-1].file_id)
    await file.download_to_drive(file_path)

    await upload(update, context, file_path)


def extract_command_arguments(message: str):
    message = message.replace(BOT_NAME, '')
    parts = message.split(' ')
    if len(parts) > 1:
        return parts[0], parts[1:]
    return parts[0], []


async def handle_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    message = update.message.text
    command, args = extract_command_arguments(message)
    context.args = args

    if command.startswith('/frame'):
        await frame(update, context)
    elif command.startswith('/random'):
        await random(update, context)
    elif command.startswith('/start'):
        await start(update, context)
    else:
        await handle_smart_reply(update, context)


def start_bot(config_path: str):
    asyncio.set_event_loop(asyncio.new_event_loop())
    config = configparser.ConfigParser()
    config.read(config_path)
    global BOT_NAME, API_URL
    BOT_NAME = config['Bot']['name']
    API_URL = config['Bot']['api_url']
    application = ApplicationBuilder().token(config['Telegram Bot API']['token']).build()
    
    application.add_handler(CommandHandler('help', help_command))
    application.add_handler(CommandHandler('start', start))
    application.add_handler(CommandHandler('frame', frame))
    application.add_handler(CommandHandler('random', random))
    application.add_handler(MessageHandler(filters.Document.Category('image/'), image_file_downloader))
    application.add_handler(MessageHandler(filters.PHOTO, image_downloader))
    application.add_handler(MessageHandler(filters.ChatType.GROUPS, handle_command))
    application.add_handler(MessageHandler(filters.ChatType.PRIVATE, handle_smart_reply))
    application.run_polling()


if __name__ == '__main__':
    import sys

    logging.basicConfig(
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        level=logging.INFO
    )
    
    start_bot(sys.argv[1])

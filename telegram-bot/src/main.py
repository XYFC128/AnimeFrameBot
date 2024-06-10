import configparser
from telegram import Update
from telegram.ext import ApplicationBuilder, ContextTypes, CommandHandler
import asyncio
from typing import Optional
from threading import Event

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(chat_id=update.effective_chat.id, text="I'm a bot, please talk to me!")


def start_bot(config_path: str):
    asyncio.set_event_loop(asyncio.new_event_loop())
    config = configparser.ConfigParser()
    config.read(config_path)
    application = ApplicationBuilder().token(config['Telegram Bot API']['token']).build()
    
    start_handler = CommandHandler('start', start)
    application.add_handler(start_handler)
    application.run_polling()


if __name__ == '__main__':
    import sys
    import logging

    logging.basicConfig(
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        level=logging.INFO
    )
    start_bot(sys.argv[1])

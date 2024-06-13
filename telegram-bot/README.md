# Telegram Anime Frame Bot

## Prerequisite

Python 3.10 or later is required.

```sh
pip install -r requirements.txt
```

### Create a Bot in Telegram

Talk to [BotFather](https://telegram.me/BotFather) to create a new bot. You should use `/setprivacy` command to disable group privacy if you wish to use smart reply feature in group chat.

## Bot Configuration

Create a `config.ini` file in this folder with the following content:

```ini
[Bot]
name = your_bot_name
api_url = http://localhost:8763

[Telegram Bot API]
token = your_token

[Telegram API]
api_id = your_api_id
api_hash = your_api_hash
session_string = your_session_string
```

where

- The bot name is the username of your bot, starting with `@`. For example: `@AnimeFrameBot`.
- api_url is the url for our api-server
- Bot API token can be obtained from [BotFather](https://telegram.me/BotFather).
- Configs in the Telegram API section is optional, they are used for integration testing. API id and hash can be obtained from the **API development tools** section in [Telegram Developer Portal](https://my.telegram.org).
- The session string is the login session identifier like a cookie. It can be obtained by running `tests/get_session_string.py` and filling in your API hash and id.
- token, api_id, api_hash, session_string should be kept secret

## Running

```
python src/main.py config.ini
```

### Telegram Commands
1. `/help` - Shows help message
2. `/start` - Starts the bot
3. `/frame` {query} {N} - Gets similar N(default: 3) frames from given query
4. `/random` {N} - Gets random N(default: 3) frames
5. It can also send an image and upload it to the server, but be sure to provide a caption for the image if it's uploaded with compression.

### Running tests

Note: You should start the API server first for the integration tests to pass.

```
pytest
```

Run unit tests only and see the coverage:

```
pytest --ignore=tests/test_integration.py --cov=src
```

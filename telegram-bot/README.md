# Telegram Anime Frame Bot

## Prerequisite

```sh
pip install -r requirements.txt
```

## Bot Configuration

Create a `config.ini` file in this folder with the following content:

```ini
[Bot]
name = your_bot_name

[Telegram Bot API]
token = your_token

[Telegram API]
api_id = your_api_id
api_hash = your_api_hash
session_string = your_session_string
```

where

- The bot name is the username of your bot, starting with `@`. For example: `@AnimeFrameBot`.
- Bot API token can be obtained from [BotFather](https://telegram.me/BotFather).
- Configs in the Telegram API section is optional, they are used for integration testing. API id and hash can be obtained from the **API development tools** section in [Telegram Developer Portal](https://my.telegram.org).
- The session string is the login session identifier like a cookie. It can be obtained by running `tests/get_session_string.py` and filling in your API hash and id.
- token, api_id, api_hash, session_string should be kept secret

## Running

```
python src/main.py config.ini
```

### Running tests

```
pytest
```

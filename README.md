# Anime Frame Bot

## Introduction

### Subject

Our project is to design a telegram bot and a API server for it that can enter a keyword (words or phrases) and output a frame containing the keyword in an anime in a TDD manner. The bot server will search for the keyword in the folder and return the frame containing the keyword.

### Motivation

We are all anime fans, and we often want to find the frame of the anime containing famous lines or quotations. However, it is difficult to find the desired frame quickly. Therefore, we want to develop a program that can search the frame for us, and a Telegram bot so we can use it conveniently.

## Implementation

Our implementation contains an API server and a Telegram chatbot. The API server provides the core functionalities including searching the anime frames. The Telegram bot provides a convenient way to interact with the server using Telegram.

The server is implemented using Golang's standard net/http package and tested with Golang's built-in test utilities, including unit tests, integration tests, fuzz tests, mutation tests, and code coverage.

Then, we use Python to implement the Telegram bot and do unit tests with pytest.

At last, we use telethon to do end-to-end tests for the Telegram bot.

## File Structure
- `api-server`: The main API server go program
- `telegram-bot`: The main Telegram bot python program
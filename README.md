# Anime Frame Bot

## Introduction

### Subject

Our project is to design a telegram bot and a web frontend that can enter a keyword (words or phrases) and output a frame containing the keyword in an anime in a TDD manner. The bot server will search for the keyword in the database and return the frame containing the keyword.

### Motivation

We are all anime fans, and we often want to find the frame of the anime containing famous lines or quotations. However, it is difficult to find the desired frame quickly. Therefore, we want to develop a program that can search the frame for us, and a Telegram bot and web frontend so we can use it everywhere conveniently.

## Implementation Plan

Our implementation contains an API server, a web frontend, and a Telegram chatbot. The API server provides the core functionalities including indexing and searching the anime frames. The web frontend provides a friendly interface to view the indexed anime. The Telegram bot provides another way to interact with the server using Telegram.

The server will be implemented using Golang's standard net/http package and tested with Golang's built-in test utilities.

For the web frontend, it is simple plain HTML with CSS website. It can be tested with Selenium, Cypress, or Puppeteer.

Finally, we plan to use Python to implement the Telegram bot and test it with Pytest.

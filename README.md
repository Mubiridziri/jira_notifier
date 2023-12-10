# Jira Notifier
[![Go](https://github.com/Mubiridziri/jira_notifier/actions/workflows/go.yml/badge.svg)](https://github.com/Mubiridziri/jira_notifier/actions/workflows/go.yml)
> If you, like me, have cluttered mail from various notifications. This simple solution will allow you to receive notifications about new tasks assigned to you directly in your Telegram.

## Setup

### Creating .env file like this:
```.env
TELEGRAM_TOKEN=BotToken
DATABASE_NAME=bot.db
JIRA_ADDRESS=https://jira.example.ru
```

You can also determine the data update time in handlers using variables: `TELEGRAM_UPDATE_INTERVAL_SEC`, `JIRA_UPDATE_INTERVAL_MIN`

## Creating Telegram Bot

You can create a Telegram bot and give away a token using @BotFather

## Install dependencies 

```bash
$ go get
```

## Run dev

```bash
$ go run main.go
```

## Build application

```bash
$ go build
```

## Getting with bot

During the first dialogue with the bot, it will ask you to send a personal token. You can take it in your Jira profile settings.

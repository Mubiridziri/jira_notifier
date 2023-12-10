package models

import (
	"jiraAwesomeBot/config"
	"jiraAwesomeBot/helpers"
)

type TelegramBot struct {
	Token  string
	Buffer int
}

type TelegramUpdate struct {
	User    User
	Message Message
}

func NewTelegramBot() TelegramBot {
	return TelegramBot{
		Token:  helpers.GetEnvStr(config.TelegramToken, ""),
		Buffer: 1,
	}
}

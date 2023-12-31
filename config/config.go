package config

import "jira_notifier/helpers"

const (
	TelegramToken             = "TELEGRAM_TOKEN"
	TelegramUpdateIntervalSec = "TELEGRAM_UPDATE_INTERVAL_SEC"
	DatabaseName              = "DATABASE_NAME"
	JiraAddress               = "JIRA_ADDRESS"
	JiraUpdateInterval        = "JIRA_UPDATE_INTERVAL_MIN"
	JiraPersonalTokenPageLink = "JIRA_PERSONAL_TOKEN_PAGE_LINK"
	AdminChatId               = "ADMIN_CHAT_ID"
)

type Config struct {
	Telegram TelegramConfig
	Database DatabaseConfig
	Jira     JiraConfig
	Alive    bool
}

type TelegramConfig struct {
	Token             string
	UpdateIntervalSec int
	AdminChatId       uint
}

type DatabaseConfig struct {
	Name string
}

type JiraConfig struct {
	JiraAddress               string
	UpdateIntervalMin         int
	JiraPersonalTokenPageLink string
}

func New() *Config {
	return &Config{
		Telegram: TelegramConfig{
			Token:             helpers.GetEnvStr(TelegramToken, ""),
			UpdateIntervalSec: helpers.GetEnvInt(TelegramUpdateIntervalSec, 5),
			AdminChatId:       helpers.GetEnvUint(AdminChatId, 0),
		},
		Database: DatabaseConfig{
			Name: helpers.GetEnvStr(DatabaseName, ""),
		},
		Jira: JiraConfig{
			JiraAddress:               helpers.GetEnvStr(JiraAddress, ""),
			JiraPersonalTokenPageLink: helpers.GetEnvStr(JiraPersonalTokenPageLink, ""),
			UpdateIntervalMin:         helpers.GetEnvInt(JiraUpdateInterval, 1),
		},
		Alive: true,
	}
}

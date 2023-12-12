package config

import "jira_notifier/helpers"

const (
	TelegramToken             = "TELEGRAM_TOKEN"
	TelegramUpdateIntervalSec = "TELEGRAM_UPDATE_INTERVAL_SEC"
	DatabaseName              = "DATABASE_NAME"
	JiraAddress               = "JIRA_ADDRESS"
	JiraUpdateInterval        = "JIRA_UPDATE_INTERVAL_MIN"
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
}

type DatabaseConfig struct {
	Name string
}

type JiraConfig struct {
	JiraAddress       string
	UpdateIntervalMin int
}

func New() *Config {
	return &Config{
		Telegram: TelegramConfig{
			Token:             helpers.GetEnvStr(TelegramToken, ""),
			UpdateIntervalSec: helpers.GetEnvInt(TelegramUpdateIntervalSec, 5),
		},
		Database: DatabaseConfig{
			Name: helpers.GetEnvStr(DatabaseName, ""),
		},
		Jira: JiraConfig{
			JiraAddress:       helpers.GetEnvStr(JiraAddress, ""),
			UpdateIntervalMin: helpers.GetEnvInt(JiraUpdateInterval, 1),
		},
		Alive: true,
	}
}

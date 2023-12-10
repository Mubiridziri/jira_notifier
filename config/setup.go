package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
)

var CFG *Config

func LoadConfig() error {
	var cfg *Config
	if err := godotenv.Load(); err != nil {
		return errors.New("no .env file found")
	}

	cfg = New()
	if err := validate(cfg); err != nil {
		return err
	}
	CFG = cfg

	return nil
}

func validate(config *Config) error {
	if config.Telegram.Token == "" {
		return createEnvNotNullError(TelegramToken)
	}
	if config.Database.Name == "" {
		return createEnvNotNullError(DatabaseName)
	}
	if config.Jira.JiraAddress == "" {
		return createEnvNotNullError(JiraAddress)
	}
	return nil
}

func createEnvNotNullError(key string) error {
	return errors.New(fmt.Sprintf("env variable %v cannot be null", key))
}

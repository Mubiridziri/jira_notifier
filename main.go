package main

import (
	"fmt"
	"jiraAwesomeBot/config"
	"jiraAwesomeBot/models"
	"jiraAwesomeBot/services"
)

func main() {
	fmt.Println("Starting...")
	fmt.Println("Load .env file...")
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	fmt.Println("Connect to local database...")
	if err := models.ConnectDatabase(); err != nil {
		panic(err)
	}

	fmt.Println("Starting Telegram Bot Handler...")
	bot, err := services.StartBot()
	if err != nil {
		fmt.Println("Cannot start a bot")
		panic(err)
	}

	fmt.Println("Update database...")
	services.HandleUserIssues(bot, true)
	fmt.Println("Starting Jira Issue Listener...")
	go services.StartListener(bot, false)

	for {
		if !config.CFG.Alive {
			break
		}
	}

	fmt.Println("Shutdown...")

}

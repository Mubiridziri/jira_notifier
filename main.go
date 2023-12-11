package main

import (
	"fmt"
	"github.com/TwiN/go-color"
	"jiraAwesomeBot/config"
	"jiraAwesomeBot/models"
	"jiraAwesomeBot/services"
	"net/http"
)

func main() {
	fmt.Println(color.Ize(color.Green, "Starting..."))
	fmt.Println(color.Ize(color.Green, "Load .env file..."))
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	fmt.Println(color.Ize(color.Green, "Connect to local database..."))
	if err := models.ConnectDatabase(); err != nil {
		panic(err)
	}

	fmt.Println(color.Ize(color.Green, "Starting Telegram Bot Handler..."))
	err := services.StartBot()
	if err != nil {
		panic(err)
	}

	fmt.Println(color.Ize(color.Green, "Update database..."))
	services.PreloadUpdatesToDatabase()
	services.HandleUserIssues(true)
	fmt.Println(color.Ize(color.Green, "Starting Jira Issue Listener..."))

	//Start goroutines listener
	go services.StartBotListener()
	go services.StartJiraListener(false)

	//TODO SECURE API ENDPOINT FOR GLOBAL TELEGRAM MESSAGE FOR ACTIVE USERS
	//HTTP Server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

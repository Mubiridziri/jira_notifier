package main

import (
	"fmt"
	"github.com/TwiN/go-color"
	"jira_notifier/config"
	"jira_notifier/models"
	"jira_notifier/routers"
	"jira_notifier/services"
	"net/http"
)

func main() {
	fmt.Println(color.Ize(color.Green, "Starting..."))

	server := http.Server{
		Addr:    ":8080",
		Handler: routers.InitRouter(),
	}

	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	if err := models.ConnectDatabase(); err != nil {
		panic(err)
	}

	if err := services.StartBot(); err != nil {
		panic(err)
	}

	services.PreloadUpdatesToDatabase()
	services.HandleUserIssues(true)

	//Start goroutines listener
	go services.StartBotListener()
	go services.StartJiraListener(false)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

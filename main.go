package main

import (
	"flag"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/gin-gonic/gin"
	"jira_notifier/config"
	"jira_notifier/handlers"
	"jira_notifier/listeners"
	"jira_notifier/models"
	"jira_notifier/routers"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println(color.Ize(color.Green, "Starting..."))
	mode := flag.Bool("release-mode", false, "Enable release mode")
	port := flag.Int("port", 8080, "Port")
	flag.Parse()

	if *mode {
		gin.SetMode(gin.ReleaseMode)
	}

	server := http.Server{
		Addr:    ":" + strconv.Itoa(*port),
		Handler: routers.InitRouter(),
	}

	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	if err := models.ConnectDatabase(); err != nil {
		panic(err)
	}

	//Cache old data
	handlers.HandleMessages(true)
	handlers.HandleNewIssue(true)
	handlers.HandleWatchedIssue(true)

	//Init handlers
	listeners.StartListeners()

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

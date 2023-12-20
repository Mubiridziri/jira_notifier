package listeners

import (
	"jira_notifier/config"
	"jira_notifier/handlers"
	"time"
)

func StartListeners() {
	go startTelegramListeners()
	go startNotificationListeners()
	go startJiraNewIssuesListeners()
	go startJiraWatchedIssuesListeners()
}

func startTelegramListeners() {
	for {
		handlers.HandleMessages(false)
		time.Sleep(time.Duration(config.CFG.Telegram.UpdateIntervalSec) * time.Second)
	}
}

func startNotificationListeners() {
	for {
		handlers.HandleNotification()
		time.Sleep(time.Second)
	}
}

func startJiraNewIssuesListeners() {
	for {
		handlers.HandleNewIssue(false)
		time.Sleep(time.Duration(config.CFG.Jira.UpdateIntervalMin) * time.Minute)
	}
}

func startJiraWatchedIssuesListeners() {
	for {
		handlers.HandleWatchedIssue(false)
		time.Sleep(time.Duration(config.CFG.Jira.UpdateIntervalMin) * time.Minute)
	}
}

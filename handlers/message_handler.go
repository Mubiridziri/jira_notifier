package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jira_notifier/models"
	"jira_notifier/services/jira"
	"jira_notifier/services/telegram"
	"os"
)

func HandleMessages(silent bool) {
	if telegram.BOT == nil {
		telegram.InitBot()
	}

	messages := telegram.GetUpdates()
	newMessages := telegram.GetOnlyNewMessages(messages)
	for _, message := range newMessages {
		models.DB.Create(&message)
		if !silent {
			handleMessage(message)
		}

	}

	//for debug
	if len(newMessages) > 0 && gin.Mode() == gin.DebugMode {
		println(fmt.Sprintf("New Messages: %v", len(newMessages)))
	}
}

func handleMessage(message models.Message) {
	text := message.Text
	user := message.User

	if text == "/start" {
		createNotification(models.WelcomeMessageType, message)
		return
	}
	user, _ = models.FindUserById(user.ID)

	if !user.JiraAuthorized {
		if jira.ConnectAccount(text, &user) {
			models.DB.Save(&user)
			createNotification(models.TokenSuccessType, message)
			return
		}
		createNotification(models.WelcomeMessageType, message)
		return
	}
	handleCommands(message)
}

func handleCommands(message models.Message) {
	user := message.User

	if user.Role != models.AdminRole {
		return
	}

	switch message.Text {
	case "/super-exit":
		os.Exit(0)
	}

}

func createNotification(notificationType string, message models.Message) {
	notification := models.Notification{
		User: message.User,
		Type: notificationType,
	}
	models.DB.Create(&notification)
}

package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jira_notifier/models"
	"jira_notifier/services/telegram"
)

func HandleNotification() {
	notifications, err := models.FindAllActiveNotifications()

	if err != nil {
		panic(err)
	}

	if gin.Mode() == gin.DebugMode {
		println(fmt.Sprintf("Found %v notifications", len(notifications)))
	}

	for _, notification := range notifications {
		switch notification.Type {
		case models.WelcomeMessageType:
			telegram.SendTokenMessage(notification.User)
		case models.TokenSuccessType:
			telegram.SendTokenSuccessMessage(notification.User)
		case models.NewIssueType:
			telegram.SendAssigneeIssueMessage(notification.Issue, notification.User)
		case models.UpdatedIssueType:
			telegram.SendSimpleUpdatedIssueMessage(notification)
		case models.UpdatedCommentsIssueType:
			telegram.SendCommentsUpdatedIssueMessage(notification.Issue, notification.User)
		}
		notification.Notified = true
		models.DB.Save(&notification)
	}
}

package jira

import (
	"github.com/gin-gonic/gin"
	"jira_notifier/config"
	"jira_notifier/models"
	"jira_notifier/services/requests"
)

// ConnectAccount check token and update user fields
func ConnectAccount(token string, user *models.User) bool {
	url := getJiraProfileURL()
	headerBag := getHeaderBag(token)

	response, err := requests.MakeGetJsonRequest(url, headerBag)

	if err != nil {
		if gin.Mode() == gin.DebugMode {
			println("Error when check Jira token")
		}

		return false
	}

	user.JiraAuthorized = true
	user.JiraPersonalToken = token
	if name, exists := response["name"].(string); exists {
		user.JiraUsername = name
	}

	return true

}

func getHeaderBag(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

func getJiraProfileURL() string {
	return config.CFG.Jira.JiraAddress + "/rest/auth/latest/session"
}

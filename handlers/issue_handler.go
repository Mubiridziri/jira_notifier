package handlers

import (
	"jira_notifier/models"
	"jira_notifier/services/jira"
)

func HandleNewIssue(silent bool) {
	users := models.FindAllActiveUsers()

	for _, user := range users {
		jira.LoadNewIssues(user, silent)
	}
}

func HandleWatchedIssue(silent bool) {
	users := models.FindAllActiveUsers()

	for _, user := range users {
		jira.LoadWatchedIssues(user, silent)
	}
}

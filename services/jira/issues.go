package jira

import (
	"github.com/gin-gonic/gin"
	"jira_notifier/config"
	"jira_notifier/models"
	"jira_notifier/services/requests"
)

func LoadNewIssues(user models.User, silent bool) {
	url := getJiraIssuesURL()
	response := LoadIssues(url, "assignee in ("+user.JiraUsername+")", user)

	for _, item := range response {
		tag := item.(map[string]interface{})["key"].(string)
		fields := item.(map[string]interface{})["fields"].(map[string]interface{})
		updatedAt := fields["updated"].(string)

		issue, err := models.FindIssueByTag(tag, user)
		if err != nil {
			title := fields["summary"].(string)
			creator := fields["creator"].(map[string]interface{})["displayName"].(string)
			priority := fields["priority"].(map[string]interface{})["name"].(string)
			status := fields["status"].(map[string]interface{})["name"].(string)
			commentsCount := fields["comment"].(map[string]interface{})["total"].(float64)
			issue = models.Issue{
				Tag:          tag,
				UpdatedAt:    updatedAt,
				Title:        title,
				User:         user,
				Priority:     priority,
				Author:       creator,
				Status:       status,
				CommentCount: int(commentsCount),
			}
			models.DB.Create(&issue)
			if !silent {
				notification := models.Notification{
					Type:  models.NewIssueType,
					User:  user,
					Issue: issue,
				}
				models.DB.Create(&notification)
			}

		}

	}
}

func LoadWatchedIssues(user models.User, silent bool) {
	url := getJiraIssuesURL()
	response := LoadIssues(url, "watcher="+user.JiraUsername, user)

	for _, item := range response {
		tag := item.(map[string]interface{})["key"].(string)
		fields := item.(map[string]interface{})["fields"].(map[string]interface{})
		updatedAt := fields["updated"].(string)
		title := fields["summary"].(string)
		priority := fields["priority"].(map[string]interface{})["name"].(string)
		status := fields["status"].(map[string]interface{})["name"].(string)
		commentsCount := fields["comment"].(map[string]interface{})["total"].(float64)

		issue, err := models.FindIssueByTag(tag, user)
		if err != nil {
			creator := fields["creator"].(map[string]interface{})["displayName"].(string)
			issue = models.Issue{
				Tag:          tag,
				UpdatedAt:    updatedAt,
				Title:        title,
				User:         user,
				Priority:     priority,
				Author:       creator,
				Status:       status,
				CommentCount: int(commentsCount),
			}
			models.DB.Create(&issue)
		}

		if updatedAt != issue.UpdatedAt {
			issue.UpdatedAt = updatedAt
			issue.Title = title

			if !silent {
				notification := models.Notification{
					User:  user,
					Issue: issue,
				}

				if int(commentsCount) != issue.CommentCount {
					notification.Type = models.UpdatedCommentsIssueType
				} else {
					notification.Type = models.UpdatedIssueType
				}

				issue.Status = status
				issue.Priority = priority
				issue.CommentCount = int(commentsCount)

				models.DB.Create(&notification)
			}

			models.DB.Save(&issue)
		}

	}
}

func LoadIssues(url string, jql string, user models.User) []interface{} {
	headerBag := getHeaderBag(user.JiraPersonalToken)
	issues, err := requests.MakePostJsonRequest(url, gin.H{
		"jql":        jql,
		"maxResults": 10000,
		"fields": []string{
			"comment",
			"summary",
			"updated",
			"priority",
			"status",
			"creator",
		},
	}, headerBag)

	if err != nil {
		if gin.Mode() == gin.DebugMode {
			println("Error when requesting JIRA issues: " + url)
		}
		return []interface{}{}
	}

	return issues["issues"].([]interface{})

}

func getJiraIssuesURL() string {
	return config.CFG.Jira.JiraAddress + "/rest/api/latest/search"
}

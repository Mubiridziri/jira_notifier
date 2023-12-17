package jira

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jira_notifier/config"
	"jira_notifier/models"
	"jira_notifier/services/requests"
)

func LoadNewIssues(user models.User, silent bool) {
	url := getJiraIssuesURL(user.JiraUsername)
	response := LoadIssues(url, user)

	for _, item := range response {
		tag := item.(map[string]interface{})["key"].(string)
		fields := item.(map[string]interface{})["fields"].(map[string]interface{})
		updatedAt := fields["updated"].(string)

		issue, err := models.FindIssueByTag(tag, user)
		if err != nil {
			title := fields["summary"].(string)
			creator := fields["creator"].(map[string]interface{})["displayName"].(string)
			priority := fields["priority"].(map[string]interface{})["name"].(string)
			issue = models.Issue{
				Tag:       tag,
				UpdatedAt: updatedAt,
				Title:     title,
				User:      user,
				Priority:  priority,
				Author:    creator,
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
	url := getJiraWatchedIssuesURL(user.JiraUsername)
	response := LoadIssues(url, user)

	for _, item := range response {
		tag := item.(map[string]interface{})["key"].(string)
		fields := item.(map[string]interface{})["fields"].(map[string]interface{})
		updatedAt := fields["updated"].(string)
		title := fields["summary"].(string)
		priority := fields["priority"].(map[string]interface{})["name"].(string)

		issue, err := models.FindIssueByTag(tag, user)
		if err != nil {
			creator := fields["creator"].(map[string]interface{})["displayName"].(string)
			issue = models.Issue{
				Tag:       tag,
				UpdatedAt: updatedAt,
				Title:     title,
				User:      user,
				Priority:  priority,
				Author:    creator,
			}
			models.DB.Create(&issue)
		}

		if updatedAt != issue.UpdatedAt {
			issue.UpdatedAt = updatedAt
			issue.Priority = priority
			issue.Title = title
			models.DB.Save(&issue)

			if !silent {
				notification := models.Notification{
					Type:  models.UpdatedIssueType,
					User:  user,
					Issue: issue,
				}
				models.DB.Create(&notification)
			}
		}

	}
}

func LoadIssues(url string, user models.User) []interface{} {
	headerBag := getHeaderBag(user.JiraPersonalToken)
	issues, err := requests.MakeGetJsonRequest(url, headerBag)

	if err != nil {
		if gin.Mode() == gin.DebugMode {
			println("Error when requesting JIRA issues: " + url)
		}
		return []interface{}{}
	}

	return issues["issues"].([]interface{})

}

func getJiraIssuesURL(username string) string {
	return fmt.Sprintf(
		config.CFG.Jira.JiraAddress+"/rest/api/latest/search?jql=assignee%%20in%%20(%v)&maxResults=1000",
		username,
	)
}

func getJiraWatchedIssuesURL(username string) string {
	return fmt.Sprintf(
		config.CFG.Jira.JiraAddress+"/rest/api/latest/search?jql=watcher=%v&maxResults=1000",
		username,
	)
}

package services

import (
	"fmt"
	"jiraAwesomeBot/config"
	"jiraAwesomeBot/models"
	"time"
)

const NewIssueType = "new_issue_type"

var templates = map[string]string{
	NewIssueType: "" +
		"🆕 *Назначено*  %v %v\n\n" +
		"*Автор*: %v\n" +
		"*Приоритет*: %v\n",
}

func StartListener(bot models.TelegramBot, silentMode bool) {
	for {
		HandleUserIssues(bot, silentMode)
		time.Sleep(time.Duration(config.CFG.Jira.UpdateIntervalMin) * time.Minute)
	}
}

func HandleUserIssues(bot models.TelegramBot, silentMode bool) {
	users, err := models.FindAllActiveUsers()
	fmt.Println(fmt.Sprintf("Found %v active users", len(users)))

	if err != nil {
		panic(err)
	}

	for _, user := range users {
		newIssues := GetNewIssuesAssignedUser(user)
		if !silentMode {
			fmt.Println(fmt.Sprintf("Found %v new issues for %v", len(newIssues), user.JiraUsername))
			for _, issue := range newIssues {
				taskLink := fmt.Sprintf("%v/browse/%v", config.CFG.Jira.JiraAddress, issue.Tag)
				template := getMessageTemplateByType(NewIssueType)
				message := fmt.Sprintf(template, issue.Tag, issue.Title, issue.Author, issue.Priority)

				SendTelegramCustomMessage(bot, map[string]interface{}{
					"chat_id":    user.ChatID,
					"text":       message,
					"parse_mode": "markdown",
					"reply_markup": map[string]interface{}{
						"inline_keyboard": [][]interface{}{
							{map[string]interface{}{"text": "Открыть задачу", "url": taskLink}},
						},
					},
				})
			}
		}
	}
}

func HandleAuth(user *models.User, token string) bool {
	headers := getAuthHeaders(token)
	jiraAddress := getJiraAuthUrl()
	respMap, err := MakeGetJsonRequest(fmt.Sprintf("%v/session", jiraAddress), headers)

	if err != nil {
		fmt.Println(fmt.Sprintf("Jira API request error: %v", err.Error()))
		return false
	}

	user.JiraUsername = respMap["name"].(string)

	return true
}

func GetNewIssuesAssignedUser(user models.User) []models.NewIssue {
	allIssues := getIssuesList(user)
	newIssues := make([]models.NewIssue, 0)

	for _, item := range allIssues["issues"].([]interface{}) {
		tag := item.(map[string]interface{})["key"].(string)
		fields := item.(map[string]interface{})["fields"].(map[string]interface{})
		creator := fields["creator"].(map[string]interface{})["displayName"].(string)
		priority := fields["priority"].(map[string]interface{})["name"].(string)
		title := fields["summary"].(string)
		_, err := models.FindTaskByTag(tag)
		if err != nil {
			newIssue := models.NewIssue{
				Tag:      tag,
				Assignee: user.JiraUsername,
				Author:   creator,
				Priority: priority,
				Title:    title,
			}

			newIssues = append(newIssues, newIssue)
			newTask := newIssue.ConvertToDBIssue()
			models.DB.Create(&newTask)

		}
	}

	return newIssues
}

func getIssuesList(user models.User) map[string]interface{} {
	headers := getAuthHeaders(user.JiraPersonalToken)
	jiraAddress := getJiraUrl()
	url := fmt.Sprintf("%v/search?jql=assignee%%20in%%20(%v)&maxResults=1000", jiraAddress, user.JiraUsername)

	respMap, err := MakeGetJsonRequest(url, headers)
	if err != nil {
		fmt.Println(fmt.Sprintf("Jira API request (%v) error: %v", url, err.Error()))
		panic(err)
	}

	return respMap
}

func getAuthHeaders(token string) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", token),
	}
}

func getJiraUrl() string {
	jiraAddress := config.CFG.Jira.JiraAddress
	return fmt.Sprintf("%v/rest/api/latest", jiraAddress)
}

func getJiraAuthUrl() string {
	jiraAddress := config.CFG.Jira.JiraAddress
	return fmt.Sprintf("%v/rest/auth/latest", jiraAddress)
}

func getMessageTemplateByType(templateType string) string {
	return templates[templateType]
}

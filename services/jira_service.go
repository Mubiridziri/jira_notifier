package services

import (
	"fmt"
	"jiraAwesomeBot/config"
	"jiraAwesomeBot/models"
	"time"
)

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
				message := fmt.Sprintf("На Вас назначена новая задача %v\n%v", issue.Tag, taskLink)
				SendTelegramMessage(bot, message, user.ChatID)
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

func GetNewIssuesAssignedUser(user models.User) []models.Task {
	allIssues := getIssuesList(user)
	newIssues := make([]models.Task, 0)

	for _, item := range allIssues["issues"].([]interface{}) {
		tag := item.(map[string]interface{})["key"].(string)
		_, err := models.FindTaskByTag(tag)
		if err != nil {
			newTask := models.Task{Tag: tag, Assignee: user.JiraUsername}
			models.DB.Create(&newTask)
			newIssues = append(newIssues, newTask)
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

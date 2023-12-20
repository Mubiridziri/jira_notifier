package telegram

import (
	"encoding/json"
	"fmt"
	"jira_notifier/config"
	"jira_notifier/models"
	"os"
)

func SendSimpleUpdatedIssueMessage(issue models.Issue, user models.User) {
	filepath := getCoverForUpdatedIssue(issue)

	sendUpdatedIssueMessage(issue, user, filepath)
}

func SendCommentsUpdatedIssueMessage(issue models.Issue, user models.User) {
	filepath := getCoverForUpdatedCommentsIssue(issue)
	sendUpdatedIssueMessage(issue, user, filepath)
}

func sendUpdatedIssueMessage(issue models.Issue, user models.User, filepath string) {
	buttons := map[string]interface{}{
		"inline_keyboard": [][]interface{}{
			{map[string]interface{}{"text": "🔗 Открыть задачу", "url": issue.GetLink()}},
		},
	}

	replyMarkup, err := json.Marshal(buttons)

	if err != nil {
		panic(err)
	}

	SendTelegramCustomMessageWithImage(map[string]string{
		"chat_id":      fmt.Sprintf("%v", user.ChatID),
		"caption":      fmt.Sprintf("🔄 *Обновлено*  %v %v\n\n", issue.Tag, issue.Title),
		"parse_mode":   "markdown",
		"reply_markup": string(replyMarkup),
	}, filepath)

	os.Remove(filepath)
}

func SendAssigneeIssueMessage(issue models.Issue, user models.User) {

	filepath := getCoverForNewIssue(issue)

	buttons := map[string]interface{}{
		"inline_keyboard": [][]interface{}{
			{map[string]interface{}{"text": "🔗 Открыть задачу", "url": issue.GetLink()}},
		},
	}

	replyMarkup, err := json.Marshal(buttons)

	if err != nil {
		panic(err)
	}

	SendTelegramCustomMessageWithImage(map[string]string{
		"chat_id": fmt.Sprintf("%v", user.ChatID),
		"caption": fmt.Sprintf(""+
			"🆕 *Назначено*  %v %v\n\n"+
			"*Автор*: %v\n"+
			"*Приоритет*: %v\n", issue.Tag, issue.Title, issue.Author, issue.Priority),
		"parse_mode":   "markdown",
		"reply_markup": string(replyMarkup),
	}, filepath)

	os.Remove(filepath)
}

func SendTokenMessage(user models.User) {
	SendTelegramCustomMessage(map[string]interface{}{
		"chat_id":    user.ChatID,
		"text":       fmt.Sprintf("Привет, %v! Пришли мне персональный токер авторизации Jira.", user.Name),
		"parse_mode": "markdown",
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]interface{}{
				{map[string]interface{}{"text": "🔗 Создать токен", "url": config.CFG.Jira.JiraPersonalTokenPageLink}},
			},
		},
	})
}

func SendTokenSuccessMessage(user models.User) {
	SendTelegramPlainMessage(
		fmt.Sprintf("Авторизация за пользователем %v прошла успешно. "+
			"Начато обновление базы данных тикетов.", user.JiraUsername),
		user.ChatID,
	)
}

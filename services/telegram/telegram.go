package telegram

import (
	"fmt"
	"jira_notifier/config"
	"jira_notifier/models"
	"jira_notifier/services/requests"
)

type Bot struct {
	Token string
}

var BOT *Bot

func InitBot() {
	BOT = &Bot{
		Token: config.CFG.Telegram.Token,
	}
}

func GetUpdates() map[string]interface{} {
	url := getTelegramUrl("getUpdates")

	response, err := requests.MakeGetJsonRequest(url, make(map[string]string))

	if err != nil {
		println("Error when request telegram bot with method getUpdates")
		return make(map[string]interface{})
	}

	return response
}

func GetOnlyNewMessages(messages map[string]interface{}) []models.Message {
	newMessages := make([]models.Message, 0)

	updates := messages["result"].([]interface{})

	for _, update := range updates {
		typedUpdate := update.(map[string]interface{})
		updateID := typedUpdate["update_id"].(float64)
		messageDB, err := models.FindMessageById(uint(updateID))

		if err == nil {
			continue
		}

		message, exists := typedUpdate["message"].(map[string]interface{})

		if !exists {
			continue
		}

		//Message
		text := message["text"].(string)
		date := message["date"].(float64)

		//User
		from := message["from"].(map[string]interface{})
		name := from["first_name"].(string)
		username := from["username"].(string)
		chatID := from["id"].(float64)

		userDB, err := models.FindUserByChatId(uint(chatID))

		if err != nil {
			userDB = models.User{
				Name:     name,
				Username: username,
				ChatID:   uint(chatID),
			}
		}

		messageDB = models.Message{
			ID:   uint(updateID),
			Text: text,
			Date: uint(date),
			User: userDB,
		}

		newMessages = append(newMessages, messageDB)

	}

	return newMessages
}

func SendTelegramPlainMessage(message string, recipient uint) map[string]interface{} {
	url := getTelegramUrl("sendMessage")

	body := map[string]interface{}{
		"chat_id": recipient,
		"text":    message,
	}

	jsonMap, err := requests.MakePostJsonRequest(url, body, make(map[string]string))

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return jsonMap
}

func SendTelegramCustomMessage(body map[string]interface{}) map[string]interface{} {
	url := getTelegramUrl("sendMessage")

	jsonMap, err := requests.MakePostJsonRequest(url, body, make(map[string]string))

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return jsonMap
}

func SendTelegramCustomMessageWithImage(form map[string]string, filepath string) {
	url := getTelegramUrl("sendPhoto")
	err := requests.MakeFormData(url, form, "photo", map[string]string{
		"file1": filepath,
	})

	if err != nil {
		panic(err)
	}

	return
}

func getTelegramUrl(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%v/%v", BOT.Token, method)
}

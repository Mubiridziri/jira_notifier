package services

import (
	"errors"
	"fmt"
	"jiraAwesomeBot/config"
	"jiraAwesomeBot/models"
	"time"
)

const TelegramApiHost = "https://api.telegram.org/bot"

const TelegramGetMeMethod = "getMe"
const TelegramSendMessageMethod = "sendMessage"
const TelegramGetUpdatesMethod = "getUpdates"

func StartBot() (models.TelegramBot, error) {
	bot := models.NewTelegramBot()
	if err := validTokenCheck(bot); err != nil {
		return models.TelegramBot{}, errors.New("invalid telegramToken")
	}

	go handleUpdate(bot)

	return bot, nil
}

func getUpdatesChan(bot models.TelegramBot, ch chan models.TelegramUpdate) {
	for {
		url := getTelegramUrl(bot, TelegramGetUpdatesMethod)
		respMap, err := MakeGetJsonRequest(url, make(map[string]string))

		if err != nil {
			panic(fmt.Sprintf("Telegram API Error: %v", err.Error()))
		}

		results := respMap["result"].([]interface{})

		for _, result := range results {
			currentResult := result.(map[string]interface{})
			message, exists := currentResult["message"]
			if exists {
				from := message.(map[string]interface{})["from"].(map[string]interface{})
				text := message.(map[string]interface{})["text"].(string)
				messageId := message.(map[string]interface{})["message_id"].(float64)
				chatId := from["id"].(float64)

				_, err := models.FindMessageById(uint(messageId))

				if err != nil {
					messagedb := models.Message{
						MessageID: uint(messageId),
						Text:      text,
					}
					models.DB.Save(&messagedb)

					user, err := models.FindUserByChatId(uint(chatId))

					if err != nil {
						user.Name = from["first_name"].(string)
						user.ChatID = uint(from["id"].(float64))
						models.DB.Create(&user)
					}

					ch <- models.TelegramUpdate{
						User:    user,
						Message: messagedb,
					}
				}
			}

		}

		time.Sleep(time.Duration(config.CFG.Telegram.UpdateIntervalSec) * time.Second)
	}
}

func handleUpdate(bot models.TelegramBot) {
	ch := make(chan models.TelegramUpdate, bot.Buffer)
	go getUpdatesChan(bot, ch)

	for item := range ch {
		fmt.Println(fmt.Sprintf("Handled new messagr from %v with text %v", item.User.ChatID, item.Message.Text))
		user := item.User
		message := item.Message.Text

		if !user.JiraAuthorized {
			if message != "/start" {
				if ok := HandleAuth(&user, message); ok {
					user.JiraPersonalToken = message
					user.JiraAuthorized = true
					models.DB.Save(&user)

					_ = SendTelegramPlainMessage(
						bot,
						fmt.Sprintf("Спасибо, %v! Авторизация прошла успешно, токен сохранен.", user.Name),
						user.ChatID)
					HandleUserIssues(bot, true)
					_ = SendTelegramPlainMessage(
						bot,
						"Успешно актуализирована база данных ранее существующих тикетов.",
						user.ChatID)
					continue
				}
			}

			_ = SendTelegramPlainMessage(
				bot,
				fmt.Sprintf("Привет, %v! Пришли мне персональный токер авторизации Jira.", user.Name),
				user.ChatID)
		}
		if message == "/super-exit" {
			config.CFG.Alive = false
		}

	}
}

func validTokenCheck(bot models.TelegramBot) error {
	url := getTelegramUrl(bot, TelegramGetMeMethod)

	if _, err := MakeGetJsonRequest(url, make(map[string]string)); err != nil {
		return err
	}

	return nil
}

func SendTelegramPlainMessage(bot models.TelegramBot, message string, recipient uint) map[string]interface{} {
	url := getTelegramUrl(bot, TelegramSendMessageMethod)

	body := map[string]interface{}{
		"chat_id": recipient,
		"text":    message,
	}

	jsonMap, err := MakePostJsonRequest(url, body, make(map[string]string))

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return jsonMap
}

func SendTelegramCustomMessage(bot models.TelegramBot, body map[string]interface{}) map[string]interface{} {
	url := getTelegramUrl(bot, TelegramSendMessageMethod)

	jsonMap, err := MakePostJsonRequest(url, body, make(map[string]string))

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return jsonMap
}

func getTelegramUrl(bot models.TelegramBot, method string) string {
	return fmt.Sprintf("%v%v/%v", TelegramApiHost, bot.Token, method)
}

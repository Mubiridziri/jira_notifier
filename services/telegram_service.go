package services

import (
	"errors"
	"fmt"
	"github.com/TwiN/go-color"
	"jira_notifier/config"
	"jira_notifier/models"
	"time"
)

const TelegramApiHost = "https://api.telegram.org/bot"

const TelegramGetMeMethod = "getMe"
const TelegramSendMessageMethod = "sendMessage"
const TelegramGetUpdatesMethod = "getUpdates"

var BOT *models.TelegramBot

func StartBot() error {
	BOT = models.NewTelegramBot()
	if err := validTokenCheck(); err != nil {
		return errors.New("invalid telegramToken")
	}

	return nil
}

func getUpdatesChan(ch chan models.TelegramUpdate) {
	for {
		url := getTelegramUrl(TelegramGetUpdatesMethod)
		respMap, err := MakeGetJsonRequest(url, make(map[string]string))

		if err != nil {
			//Wait time when read: connection reset by peer
			message := fmt.Sprintf("Telegram API Error: %v", err.Error())
			fmt.Println(color.Ize(color.Red, message))
			time.Sleep(time.Duration(config.CFG.Telegram.UpdateIntervalSec) * time.Second)
			continue

		}

		results := respMap["result"].([]interface{})
		for _, result := range results {
			currentResult := result.(map[string]interface{})
			updateId := currentResult["update_id"].(float64)

			if updateIdExists := models.IsExistUpdateId(uint(updateId)); updateIdExists {
				continue
			}
			//Save update_id to database
			models.DB.Create(&models.Update{UpdateID: uint(updateId)})

			message, isMessage := currentResult["message"]
			callbackQuery, isCallbackQuery := currentResult["callback_query"]

			if isMessage {
				handleMessage(message, ch)
			}
			if isCallbackQuery {
				handleCallbackQuery(callbackQuery)
			}

		}

		time.Sleep(time.Duration(config.CFG.Telegram.UpdateIntervalSec) * time.Second)
	}
}

func handleCallbackQuery(callbackQuery interface{}) {
	tag := callbackQuery.(map[string]interface{})["data"].(string)
	from := callbackQuery.(map[string]interface{})["from"].(map[string]interface{})
	chatId := from["id"].(float64)
	models.CreateIfNotExist(tag, uint(chatId))
}

func handleMessage(message interface{}, ch chan models.TelegramUpdate) {
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
		username := from["username"].(string)

		if err != nil {
			user.Name = from["first_name"].(string)
			user.ChatID = uint(from["id"].(float64))
			user.Username = username
			models.DB.Create(&user)
		}

		//Update username when his changed
		if user.Username != username {
			user.Username = username
			models.DB.Save(&user)
		}

		ch <- models.TelegramUpdate{
			User:    user,
			Message: messagedb,
		}
	}
}

func PreloadUpdatesToDatabase() {
	fmt.Println(color.Ize(color.Green, "Update database..."))
	url := getTelegramUrl(TelegramGetUpdatesMethod)
	respMap, err := MakeGetJsonRequest(url, make(map[string]string))

	if err != nil {
		panic(fmt.Sprintf("Telegram API Error: %v", err.Error()))
	}
	results := respMap["result"].([]interface{})
	for _, result := range results {
		currentResult := result.(map[string]interface{})
		updateId := currentResult["update_id"].(float64)

		if updateIdExists := models.IsExistUpdateId(uint(updateId)); updateIdExists {
			continue
		}
		//Save update_id to database
		models.DB.Create(&models.Update{UpdateID: uint(updateId)})
	}

}

func StartBotListener() {
	fmt.Println(color.Ize(color.Green, "Starting Telegram Listener..."))
	ch := make(chan models.TelegramUpdate, BOT.Buffer)
	go getUpdatesChan(ch)

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
						fmt.Sprintf("Спасибо, %v! Авторизация прошла успешно, токен сохранен.", user.Name),
						user.ChatID)
					HandleUserIssues(true)
					_ = SendTelegramPlainMessage(
						"Успешно актуализирована база данных ранее существующих тикетов.",
						user.ChatID)
					continue
				}
			}

			_ = SendTelegramPlainMessage(
				fmt.Sprintf("Привет, %v! Пришли мне персональный токер авторизации Jira.", user.Name),
				user.ChatID)
		}

		//Handle command only when user is admin
		if message == "/super-exit" && user.Role == models.AdminRole {
			config.CFG.Alive = false
		}

	}
}

func validTokenCheck() error {
	url := getTelegramUrl(TelegramGetMeMethod)

	if _, err := MakeGetJsonRequest(url, make(map[string]string)); err != nil {
		return err
	}

	return nil
}

func SendTelegramPlainMessage(message string, recipient uint) map[string]interface{} {
	url := getTelegramUrl(TelegramSendMessageMethod)

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

func SendTelegramCustomMessage(body map[string]interface{}) map[string]interface{} {
	url := getTelegramUrl(TelegramSendMessageMethod)

	jsonMap, err := MakePostJsonRequest(url, body, make(map[string]string))

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return jsonMap
}

func getTelegramUrl(method string) string {
	return fmt.Sprintf("%v%v/%v", TelegramApiHost, BOT.Token, method)
}

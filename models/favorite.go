package models

import "database/sql"

type Favorite struct {
	ID     uint
	Tag    string
	ChatID uint
}

func CreateIfNotExist(tag string, chatId uint) {
	err := DB.
		Where("tag = @tag AND chat_id = @chat", sql.Named("tag", tag), sql.Named("chat", chatId)).
		First(&Favorite{}).
		Error
	if err != nil {
		favorite := Favorite{Tag: tag, ChatID: chatId}
		DB.Create(&favorite)
	}
}

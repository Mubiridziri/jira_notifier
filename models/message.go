package models

import (
	"database/sql"
	"time"
)

type Message struct {
	MessageID uint `gorm:"primary_key"`
	Text      string
	CreatedAt time.Time
}

func FindMessageById(id uint) (*Message, error) {
	var message Message
	err := DB.Where("message_id = @message", sql.Named("message", id)).First(&message).Error
	if err != nil {
		return &Message{}, err
	}

	return &message, nil
}

package models

import (
	"gorm.io/gorm"
)

type Message struct {
	ID     uint `gorm:"uniqueIndex"`
	Text   string
	Date   uint
	UserID uint
	User   User
}

func FindMessageById(updateId uint) (Message, error) {
	var message Message
	err := DB.Preload("User").Where(&Message{ID: updateId}).First(&message).Error
	return message, err
}

func (message *Message) BeforeCreate(tx *gorm.DB) error {
	message.UserID = message.User.ID

	return nil
}

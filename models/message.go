package models

import (
	"database/sql"
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
	//when use .First method I catch NotFoundRecord error, why?
	err := DB.Preload("User").Where("id = @id", sql.Named("id", updateId)).Find(&message).Error
	return message, err
}

func (message *Message) BeforeCreate(tx *gorm.DB) error {
	message.UserID = message.User.ID

	return nil
}

package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type User struct {
	ID                uint `gorm:"primary_key"`
	ChatID            uint `gorm:"unique"`
	Name              string
	JiraPersonalToken string
	JiraAuthorized    bool
	JiraUsername      string
}

func FindUserByChatId(chatId uint) (User, error) {
	var user User
	err := DB.Where("chat_id = @id", sql.Named("id", chatId)).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FindAllActiveUsers() ([]User, error) {
	var users []User
	err := DB.Where(User{JiraAuthorized: true}).Find(&users).Error

	if err != nil {
		return make([]User, 0), err
	}

	return users, nil
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.JiraAuthorized = false
	return
}

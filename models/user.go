package models

import (
	"database/sql"
	"gorm.io/gorm"
	"jira_notifier/config"
)

const (
	AdminRole = "admin_role"
	UserRole  = "user_role"
)

type User struct {
	ID                uint `gorm:"primary_key"`
	ChatID            uint `gorm:"unique"`
	Username          string
	Name              string
	JiraPersonalToken string
	JiraAuthorized    bool
	JiraUsername      string
	Role              string
}

type Login struct {
	Username string `json:"username" binding:"required"`
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

func FindUserByUsername(username string) (User, error) {
	var user User
	err := DB.Where("username = @username", sql.Named("username", username)).First(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FindUserById(id uint) (User, error) {
	var user User
	err := DB.Where("id = @id", sql.Named("id", id)).First(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.JiraAuthorized = false
	user.Role = UserRole
	if user.ChatID == config.CFG.Telegram.AdminChatId {
		user.Role = AdminRole
	}

	return
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if user.ChatID == config.CFG.Telegram.AdminChatId {
		user.Role = AdminRole
	}

	return
}

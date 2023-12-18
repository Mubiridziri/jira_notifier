package models

import (
	"database/sql"
	"gorm.io/gorm"
	"jira_notifier/config"
	"jira_notifier/scopes"
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
	Messages          []Message      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Issues            []Issue        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Notifications     []Notification `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type APIUser struct {
	ID             uint
	ChatID         uint
	Username       string
	Name           string
	TelegramActive bool
	JiraActive     bool
	Role           string
}

func (user *User) ConvertAPIUser() APIUser {
	return APIUser{
		ID:             user.ID,
		ChatID:         user.ChatID,
		Username:       user.Username,
		Name:           user.Name,
		TelegramActive: len(user.Username) > 0,
		JiraActive:     user.JiraAuthorized,
		Role:           user.Role,
	}
}

func ConvertUserToAPIUser(users []User) []APIUser {
	var apiUsers []APIUser
	for _, user := range users {
		apiUsers = append(apiUsers, user.ConvertAPIUser())
	}
	return apiUsers
}

type Login struct {
	Username string `json:"username" binding:"required"`
}

func FindPaginatedUsers(page, limit int) ([]User, error) {
	var users []User
	err := DB.Scopes(scopes.GetPaginationScope(page, limit)).Find(&users).Error
	return users, err
}

func GetUsersTotal() (int64, error) {
	var count int64
	err := DB.Model(&User{}).Count(&count).Error
	return count, err
}

func FindUserByChatId(chatId uint) (User, error) {
	var user User
	err := DB.Where("chat_id = @id", sql.Named("id", chatId)).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FindAllActiveUsers() []User {
	var users []User
	err := DB.Where("jira_authorized = true").Find(&users).Error

	if err != nil {
		panic(err)
	}

	return users
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

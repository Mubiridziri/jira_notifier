package models

import (
	"errors"
	"gorm.io/gorm"
	"jira_notifier/helpers"
)

const (
	NewIssueType       = "NewIssue"
	UpdatedIssueType   = "UpdatedIssue"
	TokenSuccessType   = "TokenSuccessType"
	DataLoadedType     = "DataLoaded"
	WelcomeMessageType = "WelcomeMessage"
)

type Notification struct {
	ID       uint
	UserID   uint
	User     User
	Notified bool
	Type     string
	IssueID  uint
	Issue    Issue
}

func FindAllActiveNotifications() ([]Notification, error) {
	var notifications []Notification
	err := DB.Preload("User").Preload("Issue").Preload("Issue.User").Where("notified = false").Find(&notifications).Error
	return notifications, err
}

func (notification *Notification) BeforeCreate(tx *gorm.DB) error {
	if helpers.InStringArray(
		[]string{NewIssueType, UpdatedIssueType, TokenSuccessType, WelcomeMessageType, DataLoadedType},
		notification.Type) == -1 {
		return errors.New("invalid notification type")
	}
	notification.UserID = notification.User.ID
	notification.IssueID = notification.Issue.ID
	notification.Notified = false

	return nil
}

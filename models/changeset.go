package models

import "gorm.io/gorm"

type ChangeSet struct {
	ID             uint
	IssueID        uint
	Issue          Issue
	Field          string
	OldValue       string
	NewValue       string
	NotificationID uint
	Notification   Notification
}

func NewChangeSet(issue Issue, field, oldValue, newValue string, notification Notification) ChangeSet {
	return ChangeSet{
		Issue:        issue,
		Notification: notification,
		Field:        field,
		OldValue:     oldValue,
		NewValue:     newValue,
	}
}

func (changeSet *ChangeSet) BeforeCreate(tx *gorm.DB) error {
	changeSet.NotificationID = changeSet.Notification.ID
	changeSet.IssueID = changeSet.Issue.ID

	return nil
}

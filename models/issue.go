package models

import (
	"database/sql"
	"gorm.io/gorm"
	"jira_notifier/config"
)

type Issue struct {
	ID            uint `gorm:"primary_key"`
	Tag           string
	Title         string
	UpdatedAt     string
	Priority      string
	Author        string
	CommentCount  int
	Status        string
	UserID        uint
	User          User
	Notifications []Notification `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ChangeSets    []ChangeSet    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func FindIssueByTag(tag string, user User) (Issue, error) {
	var issue Issue
	err := DB.Preload("User").Where("tag = @tag AND user_id = @user",
		sql.Named("tag", tag), sql.Named("user", user.ID)).First(&issue).Error
	if err != nil {
		return Issue{}, err
	}

	return issue, nil
}

func (issue *Issue) GetLink() string {
	return config.CFG.Jira.JiraAddress + "/browse/" + issue.Tag
}

func (issue *Issue) BeforeCreate(tx *gorm.DB) error {
	issue.UserID = issue.User.ID
	return nil
}

package models

import (
	"database/sql"
)

type TrackedIssue struct {
	ID        uint `gorm:"primary_key"`
	Tag       string
	Username  string
	UpdatedAt string
}

func FindTrackedIssueByTag(tag, username string) (*TrackedIssue, error) {
	var task *TrackedIssue
	if err := DB.
		Where(
			"tag = @tag AND username = @username",
			sql.Named("tag", tag),
			sql.Named("username", username)).
		First(&task).Error; err != nil {
		return &TrackedIssue{}, err
	}
	return task, nil
}

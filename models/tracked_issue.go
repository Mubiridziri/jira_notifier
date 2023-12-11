package models

import (
	"database/sql"
)

type TrackedIssue struct {
	ID        uint   `gorm:"primary_key"`
	Tag       string `gorm:"unique"`
	UpdatedAt string
}

func FindTrackedIssueByTag(tag string) (*TrackedIssue, error) {
	var task *TrackedIssue
	if err := DB.Where("tag = @tag", sql.Named("tag", tag)).First(&task).Error; err != nil {
		return &TrackedIssue{}, err
	}
	return task, nil
}

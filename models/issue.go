package models

import "database/sql"

type Issue struct {
	ID       uint   `gorm:"primary_key"`
	Tag      string `gorm:"unique"`
	Assignee string
}

type NewIssue struct {
	Tag      string
	Title    string
	Author   string
	Priority string
	Assignee string
}

func (issue *NewIssue) ConvertToDBIssue() Issue {
	return Issue{
		Tag:      issue.Tag,
		Assignee: issue.Assignee,
	}
}

func FindTaskByTag(tag string) (*Issue, error) {
	var task *Issue
	if err := DB.Where("tag = @tag", sql.Named("tag", tag)).First(&task).Error; err != nil {
		return &Issue{}, err
	}
	return task, nil
}

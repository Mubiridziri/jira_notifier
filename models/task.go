package models

import "database/sql"

type Task struct {
	ID       uint   `gorm:"primary_key"`
	Tag      string `gorm:"unique"`
	Assignee string
}

func FindTaskByTag(tag string) (*Task, error) {
	var task *Task
	if err := DB.Where("tag = @tag", sql.Named("tag", tag)).First(&task).Error; err != nil {
		return &Task{}, err
	}
	return task, nil
}

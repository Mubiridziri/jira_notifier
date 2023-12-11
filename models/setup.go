package models

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"jiraAwesomeBot/config"
)

var DB *gorm.DB

func ConnectDatabase() error {
	database, err := gorm.Open(sqlite.Open(config.CFG.Database.Name), &gorm.Config{})

	if err != nil {
		return errors.New("failed connect to database")
	}

	err = database.AutoMigrate(&User{})
	err = database.AutoMigrate(&Message{})
	err = database.AutoMigrate(&Issue{})
	err = database.AutoMigrate(&Update{})
	err = database.AutoMigrate(&Favorite{})
	err = database.AutoMigrate(&TrackedIssue{})

	if err != nil {
		return errors.New("failed auto migrate database")
	}

	DB = database

	return nil
}

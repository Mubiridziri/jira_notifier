package models

import (
	"errors"
	"fmt"
	"github.com/TwiN/go-color"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"jira_notifier/config"
)

var DB *gorm.DB

func ConnectDatabase() error {
	fmt.Println(color.Ize(color.Green, "Connect to local database..."))
	database, err := gorm.Open(sqlite.Open(config.CFG.Database.Name), &gorm.Config{})

	if err != nil {
		return errors.New("failed connect to database")
	}

	err = database.AutoMigrate(&User{})
	err = database.AutoMigrate(&Message{})
	err = database.AutoMigrate(&Issue{})
	err = database.AutoMigrate(&Update{})
	err = database.AutoMigrate(&TrackedIssue{})
	err = database.AutoMigrate(&AuthorizationCode{})

	//TODO AFTER RELEASE 0.4
	migrator := database.Migrator()
	if migrator.HasTable("favorites") {
		err = migrator.DropTable("favorites")
	}

	if err != nil {
		return errors.New("failed auto migrate database")
	}

	DB = database

	return nil
}

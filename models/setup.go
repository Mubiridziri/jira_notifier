package models

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"jira_notifier/config"
)

var DB *gorm.DB

func ConnectDatabase() error {
	database, err := gorm.Open(sqlite.Open(config.CFG.Database.Name), &gorm.Config{})

	if err != nil {
		return errors.New("failed connect to database")
	}

	//TODO REMOVE AFTER RELEASE 0.4
	migrator := database.Migrator()
	if migrator.HasTable("favorites") {
		err = migrator.DropTable("favorites")
	}
	if migrator.HasTable("issues") {
		err = migrator.DropTable("issues")
	}
	if migrator.HasTable("authorization_codes") {
		err = migrator.DropTable("authorization_codes")
	}
	if migrator.HasTable("messages") {
		err = migrator.DropTable("messages")
	}
	if migrator.HasTable("tracked_issues") {
		err = migrator.DropTable("tracked_issues")
	}
	if migrator.HasTable("updates") {
		err = migrator.DropTable("updates")
	}

	err = database.AutoMigrate(&User{})
	err = database.AutoMigrate(&UserCode{})
	err = database.AutoMigrate(&Issue{})
	err = database.AutoMigrate(&Message{})
	err = database.AutoMigrate(&Notification{})

	if err != nil {
		return errors.New("failed auto migrate database")
	}

	DB = database

	return nil
}

package models

import (
	"database/sql"
	"fmt"
)

type Update struct {
	ID       uint
	UpdateID uint
}

func IsExistUpdateId(updateId uint) bool {
	if err := DB.Where("update_id = @id", sql.Named("id", updateId)).First(&Update{}).Error; err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

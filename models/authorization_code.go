package models

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type AuthorizationCode struct {
	Code      int
	ExpiredAt time.Time
	UserID    uint `gorm:"not null"`
}

type LoginCode struct {
	Code int `json:"code" binding:"required"`
}

func FindValidCode(code int) (AuthorizationCode, error) {
	var authorizationCode AuthorizationCode
	now := time.Now()
	err := DB.
		Where(
			"code = @code AND expired_at > @now",
			sql.Named("code", code),
			sql.Named("now", now)).
		First(&authorizationCode).
		Error

	if err != nil {
		return AuthorizationCode{}, err
	}

	return authorizationCode, nil
}

func (code *AuthorizationCode) BeforeCreate(tx *gorm.DB) error {
	code.ExpiredAt = time.Now().Add(time.Minute * 2) //add 2 minutes
	return nil
}

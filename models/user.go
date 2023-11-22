package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
	Name     string    `gorm:"type:varchar(255);not null"`
	Email    string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"`

	OTPEnabled  bool `gorm:"default:false"`
	OTPVerified bool `gorm:"default:false"`

	OTPSecret  string
	OTPAuthUrl string
}

func (user *User) BeforeCreate(*gorm.DB) error {
	user.ID = uuid.NewV4()

	return nil
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPInput struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

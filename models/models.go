package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type AuthSession struct {
	gorm.Model
	UserID        uint
	User          User
	RefreshToken  string
	ExpirationUtc time.Time
}

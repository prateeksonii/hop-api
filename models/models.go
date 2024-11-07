package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null;unique"`
	Username string `gorm:"not null;unique"`
	Password string `json:"-"`
}

type AuthSession struct {
	gorm.Model
	UserID        uint `gorm:"not null"`
	User          User
	RefreshToken  string    `gorm:"not null"`
	ExpirationUtc time.Time `gorm:"not null"`
}

type Contact struct {
	gorm.Model
	Name          *string
	UserID        uint `gorm:"not null"`
	User          User
	ContactUserID uint `gorm:"not null"`
	ContactUser   User
}

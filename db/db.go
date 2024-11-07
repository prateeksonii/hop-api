package db

import (
	"drop/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db

	migrate()
}

func migrate() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.AuthSession{})
	DB.AutoMigrate(&models.Contact{})
}

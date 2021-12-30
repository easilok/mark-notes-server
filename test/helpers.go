package test

import (
	"os"

	"github.com/easilok/mark-notes-server/models"
	"github.com/jinzhu/gorm"
)

const TEST_DATABASE_PATH string = "notes_test.db"

func CreateExampleUser(db *gorm.DB) {
	var user models.User
	if err := db.First(&user).Error; err != nil {
		user.Email = "test@test.com"
		user.Name = "test"
		user.Password = "123456"
		db.Save(&user)
	}
}

func ConnectTestDatabase() *gorm.DB {
	dbPath := TEST_DATABASE_PATH
	database, err := gorm.Open("sqlite3", dbPath)

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&models.NoteInformation{})
	database.AutoMigrate(&models.Category{})
	database.AutoMigrate(&models.User{})

	return database
}

func RemoveTestDatabase() {
	os.Remove(TEST_DATABASE_PATH)
}

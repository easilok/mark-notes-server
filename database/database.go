package database

import (
	"github.com/easilok/mark-notes-server/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

func FirstSetup(db *gorm.DB) {
	var user models.User
	hashedPassword, err := models.HashPassword("123456")
  if err != nil {
    return;
  }
	if err := db.First(&user).Error; err != nil {
		user.Email = "test@test.com"
		user.Name = "test"
		user.Password = hashedPassword
		db.Save(&user)
	}
}

func ConnectDatabase() *gorm.DB {
	dbPath := "notes" + string(os.PathSeparator) + "notes.db"
	database, err := gorm.Open("sqlite3", dbPath)

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&models.NoteInformation{})
	database.AutoMigrate(&models.Category{})
	database.AutoMigrate(&models.User{})

	FirstSetup(database)

	return database
}

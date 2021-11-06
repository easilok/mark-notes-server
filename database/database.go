package database

import (
    "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/easilok/mark-notes-server/models"
)

func FirstSetup (db *gorm.DB) {
  var user models.User
  if err := db.First(&user).Error; err != nil {
    user.Email = "test@test.com"
    user.Name = "test"
    user.Password = "123456"
    db.Save(&user)
  }
}

func ConnectDatabase() *gorm.DB {
  database, err := gorm.Open("sqlite3", "notes.db")

  if err != nil {
    panic("Failed to connect to database!")
  }

  database.AutoMigrate(&models.NoteInformation{})
  database.AutoMigrate(&models.Category{})
  database.AutoMigrate(&models.User{})
  
  FirstSetup(database)

  return database
}

package database

import (
    "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/easilok/mark-notes-server/models"
)

func ConnectDatabase() *gorm.DB {
  database, err := gorm.Open("sqlite3", "notes.db")

  if err != nil {
    panic("Failed to connect to database!")
  }

  database.AutoMigrate(&models.NoteInformation{})
  database.AutoMigrate(&models.Category{})

  return database
}

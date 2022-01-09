package test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/easilok/mark-notes-server/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const TEST_DATABASE_PATH string = "notes_test.db"

func CreateExampleUser(db *gorm.DB) models.User {
	var user models.User
	if err := db.First(&user).Error; err != nil {
		user.Email = "test@test.com"
		user.Name = "test"
		hashedPassword, _ := models.HashPassword("123456")
		user.Password = hashedPassword
		db.Save(&user)
	}
	return user
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

func MockJsonPost(c *gin.Context /* the test context */, content interface{}, method string) error {
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		return err
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

	return nil
}

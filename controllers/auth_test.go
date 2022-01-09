package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/easilok/mark-notes-server/test"
	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	db := test.ConnectTestDatabase()
	defer test.RemoveTestDatabase()

	controllers := NewBaseHandler(db)

	createdTestUser := test.CreateExampleUser(db)

	t.Run("login failed with no data", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		controllers.Login(c)

		if w.Code != http.StatusUnprocessableEntity {
			t.Fatalf("Expected status code %v got %v", http.StatusUnprocessableEntity, w.Code)
		}

		expectedMessage := "Invalid user provided"
		if !strings.Contains(w.Body.String(), expectedMessage) {
			t.Fatalf("Expected Body of %v got %v", expectedMessage, w.Body)
		}

	})

	t.Run("login failed with wrong user", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		testUser := User{
			Username: "test2@test.com",
			Password: "123456",
		}
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		err := test.MockJsonPost(c, testUser, "POST")
		if err != nil {
			t.Fatalf("Error marshaling test user")
		}

		controllers.Login(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code %v got %v", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("login succeed with right user wrong password", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		testUser := User{
			Username: createdTestUser.Email,
			Password: "1234567",
		}
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		err := test.MockJsonPost(c, testUser, "POST")
		if err != nil {
			t.Fatalf("Error marshaling test user")
		}

		controllers.Login(c)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status code %v got %v", http.StatusUnauthorized, w.Code)
			// t.Fatalf("Expected status code %v got %v, with body: %v", http.StatusUnauthorized, w.Code, w.Body)
		}

		expectedMessage := "Please provide valid login details"
		if !strings.Contains(w.Body.String(), expectedMessage) {
			t.Fatalf("Expected Body of %v got %v", expectedMessage, w.Body)
		}
	})

	t.Run("login succeed with right user", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		createdTestUser := test.CreateExampleUser(db)
		testUser := User{
			Username: createdTestUser.Email,
			Password: "123456",
		}
		c.Request = &http.Request{
			Header: make(http.Header),
		}
		err := test.MockJsonPost(c, testUser, "POST")
		if err != nil {
			t.Fatalf("Error marshaling test user")
		}

		controllers.Login(c)

		if w.Code != http.StatusOK {
			// t.Fatalf("Expected status code %v got %v, with body: %v", http.StatusOK, w.Code, w.Body)
			t.Fatalf("Expected status code %v got %v", http.StatusOK, w.Code)
		}

		expectedMessage := "access_token"
		if !strings.Contains(w.Body.String(), expectedMessage) {
			t.Fatalf("Expected Body of %v got %v", expectedMessage, w.Body)
		}
		expectedMessage = "refresh_token"
		if !strings.Contains(w.Body.String(), expectedMessage) {
			t.Fatalf("Expected Body of %v got %v", expectedMessage, w.Body)
		}
	})
}

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	c "github.com/easilok/mark-notes-server/controllers"
	"github.com/easilok/mark-notes-server/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// func TestGetCatalogForbidden(t *testing.T) {

// 	w := httptest.NewRecorder()
// 	gin.SetMode(gin.TestMode)
// 	c, _ := gin.CreateTestContext(w)
// 	HelloWorld(c)

// 	t.Run("get json data", func(t *testing.T) {
// 		assert.Equal(t, 200, w.Code)
// 	})
// }

func TestGetCatalogNothing(t *testing.T) {
	db := test.ConnectTestDatabase()
	defer test.RemoveTestDatabase()

	w := httptest.NewRecorder()
	r := gin.Default()
	gin.SetMode(gin.TestMode)

	controllers := c.NewBaseHandler(db)

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/catalog", TokenAuthMiddleware(), controllers.GetNotes)
	}

	t.Run("get json data", func(t *testing.T) {
		assert.Equal(t, 200, w.Code)
	})
}

func TestRouter(t *testing.T) {
	db := test.ConnectTestDatabase()
	defer test.RemoveTestDatabase()

	// w := httptest.NewRecorder()
	r := gin.Default()
	gin.SetMode(gin.TestMode)

	controllers := c.NewBaseHandler(db)

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/catalog", TokenAuthMiddleware(), controllers.GetNotes)
	}
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("get catalog without login", func(t *testing.T) {
		// Make a request to our server with the {base url}/ping
		resp, err := http.Get(fmt.Sprintf("%s/ping", ts.URL))

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.StatusCode != 404 {
			t.Fatalf("Expected status code 404, got %v", resp.StatusCode)
		}
	})
	// t.Run("example with full verification", func(t *testing.T) {
	// 	// Make a request to our server with the {base url}/ping
	// 	resp, err := http.Get(fmt.Sprintf("%s/ping", ts.URL))

	// 	if err != nil {
	// 		t.Fatalf("Expected no error, got %v", err)
	// 	}

	// 	if resp.StatusCode != 200 {
	// 		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	// 	}

	// 	val, ok := resp.Header["Content-Type"]

	// 	// Assert that the "content-type" header is actually set
	// 	if !ok {
	// 		t.Fatalf("Expected Content-Type header to be set")
	// 	}

	// 	// Assert that it was set as expected
	// 	if val[0] != "application/json; charset=utf-8" {
	// 		t.Fatalf("Expected \"application/json; charset=utf-8\", got %s", val[0])
	// 	}
	// })
	// t.Run("get json data", func(t *testing.T) {
	// 	// Make a request to our server with the {base url}/ping
	// 	http.Get(fmt.Sprintf("%s/ping", ts.URL))

	// 	assert.Equal(t, 200, w.Code)
	// })
}

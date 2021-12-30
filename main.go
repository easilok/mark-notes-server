package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	c "github.com/easilok/mark-notes-server/controllers"
	"github.com/easilok/mark-notes-server/database"
	"github.com/easilok/mark-notes-server/helpers"
)

func main() {
	fmt.Println("Loading dotenv")
	godotenv.Load()

	r := gin.Default()
	r.Use(CORSMiddleware())

	db := database.ConnectDatabase()
	controllers := c.NewBaseHandler(db)

	database.Initialize()

	helpers.CheckTokenSecrets()

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/catalog", TokenAuthMiddleware(), controllers.GetNotes)
		apiGroup.PATCH("/favorites/:filename", TokenAuthMiddleware(), controllers.FavoriteNote)
		apiGroup.GET("/note/:filename", TokenAuthMiddleware(), controllers.GetNote)
		apiGroup.PUT("/note/:filename", TokenAuthMiddleware(), controllers.UpdateNote)
		apiGroup.DELETE("/note/:filename", TokenAuthMiddleware(), controllers.DeleteNote)
		apiGroup.GET("/note/scan", TokenAuthMiddleware(), controllers.ScanNotes)
	}

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", controllers.Login)
		authGroup.GET("/logout", TokenAuthMiddleware(), controllers.Logout)
		authGroup.POST("/refresh", controllers.Refresh)
		authGroup.POST("/register", TokenAuthMiddleware(), controllers.Register)
	}

	r.Run("0.0.0.0:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		au, err := helpers.ExtractTokenMetadata(c.Request)
		// err := helpers.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("userId", au.UserId)
		c.Next()
	}
}

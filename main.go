package main

import (
   "github.com/gin-gonic/gin"

  "github.com/easilok/mark-notes-server/database"
  c "github.com/easilok/mark-notes-server/controllers"

)

func main() {
  r := gin.Default()

  db := database.ConnectDatabase()
  controllers := c.NewBaseHandler(db)

  apiGroup := r.Group("/api")
  {
    apiGroup.GET("/catalog", controllers.GetNotes)
    // r.GET("/books/:id", controllers.FindBook)
    // r.POST("/books", controllers.CreateBook)
    // r.PATCH("/books/:id", controllers.UpdateBook)
    // r.DELETE("/books/:id", controllers.DeleteBook)
  }

  r.Run("0.0.0.0:8080")
}


package main

import (
   "github.com/gin-gonic/gin"

  "github.com/easilok/mark-notes-server/database"
  c "github.com/easilok/mark-notes-server/controllers"

)

func main() {
  r := gin.Default()
  r.Use(CORSMiddleware())

  db := database.ConnectDatabase()
  controllers := c.NewBaseHandler(db)

  apiGroup := r.Group("/api")
  {
    apiGroup.GET("/catalog", controllers.GetNotes)
    apiGroup.PATCH("/favorites/:filename", controllers.FavoriteNote)
    apiGroup.GET("/note/:filename", controllers.GetNote)
    apiGroup.PUT("/note/:filename", controllers.UpdateNote)
    apiGroup.DELETE("/note/:filename", controllers.DeleteNote)
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

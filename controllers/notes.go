package controllers

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  "github.com/easilok/mark-notes-server/models"
)

// BaseHandler will hold everything that controller needs
type BaseHandler struct {
  db *gorm.DB
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *gorm.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

type Catalog struct {
  Notes []models.NoteInformation `json:"notes"`
  Categories []models.Category `json:"categories"`
}

type UpdateNoteInput struct {
  // Filename   string  `json:"filename" binding:"required"` // this is the param?
  Content  string  `json:"content" binding:"required"`
}


// GET /catalog
// Get notes and categories catalog
func (h *BaseHandler) GetNotes(c *gin.Context) {
  var catalog Catalog
  h.db.Find(&catalog.Notes)
  h.db.Find(&catalog.Categories)


  c.JSON(http.StatusOK, gin.H{"data": catalog})
}

// GET /note/:filename
// Get a note
func (h *BaseHandler) FindBook(c *gin.Context) {  // Get model if exist

  // Find filename on local machine
  // filename := c.Param("filename")

  // If filename not found delete from note information

  // Return note content

  c.JSON(http.StatusOK, gin.H{"data": "dummy"})
}

// PATCH /note/:filename
// Update a note
func (h *BaseHandler) UpdateBook(c *gin.Context) {

  // Validate input
  var input UpdateNoteInput
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  // Find filename on local machine
  // filename := c.Param("filename")
  
  // if filemane exists on storage -> update content -> update note information title

  // if filename does not exist on storage -> create it -> append to note information

  c.JSON(http.StatusOK, gin.H{"data": "dummy"})
}

// DELETE /note/:filename
// Delete a note
func (h * BaseHandler) DeleteBook(c *gin.Context) {
  // Find filename on local machine
  // filename := c.Param("filename")

  // if filename exists on storage -> delete it -> remove from note information

  c.JSON(http.StatusOK, gin.H{"data": true})
}

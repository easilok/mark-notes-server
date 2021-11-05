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

type UpdateNoteInput struct {
  // Filename   string  `json:"filename" binding:"required"` // this is the param?
  Content  string  `json:"content" binding:"required"`
}

type NoteInformationAPI struct {
  // ID        uint    `json:"id"`
  Filename  string  `json:"filename"`
  Title     string  `json:"tile"`
  Favorite  bool    `json:"favorite"`
}

type Catalog struct {
  Notes []NoteInformationAPI `json:"notes"`
  Categories []models.Category `json:"categories"`
}

// GET /catalog
// Get notes and categories catalog
func (h *BaseHandler) GetNotes(c *gin.Context) {
  var catalog Catalog
  h.db.Model(&models.NoteInformation{}).Scan(&catalog.Notes)
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
  filename := c.Param("filename")
  
  var editingNote models.NoteInformation
  isNewNote := false
  if err := h.db.Where("filename = ?", filename).First(&editingNote).Error; err != nil {
    isNewNote = true
  }

  if isNewNote {
    // Let's add it
    // TODO - Add note file
    editingNote.Filename = filename
    // TODO - Extract Title from content
    editingNote.Title = "Test"
    editingNote.Favorite = false
    editingNote.UserID = 1
    h.db.Create(&editingNote)
    c.JSON(http.StatusOK, gin.H{"data": editingNote})
  } else {
    // This is an update
    // TODO - Update content in the note File 
    // TODO - Extract Title from content
    editingNote.Title = "Updated Test"
    h.db.Model(&editingNote).Updates(editingNote)
    c.JSON(http.StatusOK, gin.H{"data": "note exists!"})
  }

}

// DELETE /note/:filename
// Delete a note
func (h * BaseHandler) DeleteBook(c *gin.Context) {
  // Find filename on local machine
  filename := c.Param("filename")

  // if filename exists on storage -> delete it -> remove from note information
  var deletingNote models.NoteInformation
  if err := h.db.Where("filename = ?", filename).First(&deletingNote).Error; err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
    return
  }

  h.db.Delete(&deletingNote)
  // TODO - Delete note from filesystem

  c.JSON(http.StatusOK, gin.H{"data": true})
}

package controllers

import (
  "net/http"
  "os"
  str "strings"
  // "io/ioutil"

  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  "github.com/easilok/mark-notes-server/models"
  "github.com/easilok/mark-notes-server/helpers"
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

type FavoritesInput struct {
  Favorite  *bool  `json:"favorite" binding:"required"`
}

type Catalog struct {
  Notes []models.NoteInformationAPI `json:"notes"`
  Categories []models.Category `json:"categories"`
}

// GET /catalog
// Get notes and categories catalog
func (h *BaseHandler) GetNotes(c *gin.Context) {
  var catalog Catalog
  h.db.Model(&models.NoteInformation{}).Scan(&catalog.Notes)

  // Build a response object, to allow more data
  var response = map[string]interface{}{ 
    "notes": catalog.Notes,
  }
  c.JSON(http.StatusOK, gin.H{"data": response})
}

// Patch /favorites/:filename
// Update a favorite status on note
func (h *BaseHandler) FavoriteNote(c *gin.Context) {
  // Validate input
  var input FavoritesInput
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  filename := c.Param("filename")
  filename = helpers.String(filename).GetFilename(".md")

  var selectedNote models.NoteInformation
  if err := h.db.Where("filename = ?", filename).First(&selectedNote).Error; err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
    return
  }

  selectedNote.Favorite = *input.Favorite
  // This update with map is for avoid error on update with struct working only for non zero values
  if err := h.db.Model(&selectedNote).Updates(map[string]bool{"favorite": selectedNote.Favorite}).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }

  c.JSON(http.StatusOK, gin.H{"data": selectedNote.ExportedFields()})
}

// GET /note/:filename
// Get a note
func (h *BaseHandler) GetNote(c *gin.Context) {  // Get model if exist

  // Find filename on local machine
  filename := c.Param("filename")
  filename = helpers.String(filename).GetFilename(".md")
  filepath := "notes" + string(os.PathSeparator) + filename + ".md"

  // If filename not found delete from note information
  data, err := os.ReadFile(filepath)
  if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    return
  }
  // Return note content
  noteContent := string(data)
  var response = map[string]string{ 
    "filename": filename, 
    "content": noteContent, 
  }

  var requestedNote models.NoteInformation
  if err := h.db.Where("filename = ?", filename).First(&requestedNote).Error; err != nil {
    // Note is note on catalog, so add it
    requestedNote.Filename = filename
    requestedNote.Title = helpers.String(noteContent).TitleFromMarkdown()
    requestedNote.Favorite = false
    requestedNote.UserID = 1
    h.db.Create(&requestedNote)
  }

  c.JSON(http.StatusOK, gin.H{"data": response})
}

// PUT /note/:filename
// Create/Update a note
func (h *BaseHandler) UpdateNote(c *gin.Context) {

  // Validate input
  var input UpdateNoteInput
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  // Find filename on local machine
  filename := c.Param("filename")
  filename = helpers.String(filename).GetFilename(".md")
  filepath := "notes" + string(os.PathSeparator) + filename + ".md"
  
  var editingNote models.NoteInformation
  var err error
  isNewNote := false
  if err = h.db.Where("filename = ?", filename).First(&editingNote).Error; err != nil {
    isNewNote = true
  }

  responseCode := http.StatusOK
  editingNote.Title = helpers.String(input.Content).TitleFromMarkdown()

  if isNewNote {
    // Let's add it
    editingNote.Filename = filename
    editingNote.Favorite = false
    editingNote.UserID = 1
    err = h.db.Create(&editingNote).Error
    responseCode = http.StatusCreated
  } else {
    // This is an update
    err = h.db.Model(&editingNote).Updates(editingNote).Error
  }

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }

  // Add note file system. If exists just override it
  // err = os.WriteFile(filepath, []byte(input.Content), 0666)
  // err = ioutil.WriteFile(filepath, []byte(input.Content), 0666)
  f, err := os.Create(filepath)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }
  defer f.Close()
  err = f.Chmod(0777)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }
  _, err = f.WriteString(input.Content)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }

  c.JSON(responseCode, gin.H{"data": editingNote.ExportedFields()})
}

// DELETE /note/:filename
// Delete a note
func (h * BaseHandler) DeleteNote(c *gin.Context) {
  // Find filename on local machine
  filename := c.Param("filename")
  filename = helpers.String(filename).GetFilename(".md")
  // filepath := "notes" + string(os.PathSeparator) + filename + ".md"

  // if filename exists on storage -> delete it -> remove from note information
  var deletingNote models.NoteInformation
  if err := h.db.Where("filename = ?", filename).First(&deletingNote).Error; err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
    return
  }

  h.db.Delete(&deletingNote)

  message := ""
  if err := helpers.TrashFile("notes", filename + ".md"); err != nil {
    message = "Error deleting file: " + err.Error()
  }

  c.JSON(http.StatusOK, gin.H{"data": true, "message": message})
}

// GET /note/scan
// Scan not indexed notes on local drive
func (h * BaseHandler) ScanNotes(c *gin.Context) {
  var registeredNotes []string;
  // Fetch all filenames in the database catalog
  h.db.Model(&models.NoteInformation{}).Pluck("filename", &registeredNotes)

  // Fecth all filenames in the notes folder
  f, err := os.Open("notes")
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  localFiles, err := f.ReadDir(0)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  var missingNotes []models.NoteInformationAPI
  // Check each file in the system if is already cataloged
  for _, v := range localFiles {
    filepath := v.Name()
    if !v.IsDir() && str.HasSuffix(filepath, ".md") {
      filename := helpers.String(filepath).GetFilename(".md")
      fileInCatalog := false
      for _, existingNote := range registeredNotes {
        if (filename == existingNote) {
          fileInCatalog = true
          break
        }
      }
      if !fileInCatalog {
        // If file is not yet cataloged, catalog it and sent back to client
        fullPath := "notes" + string(os.PathSeparator) + filepath
        data, err := os.ReadFile(fullPath)
        if err != nil {
          // if file can't be read, ignore it and try the next one
          continue
        }
        noteContent := string(data)
        var missingNote models.NoteInformation
        missingNote.Favorite = false
        missingNote.Filename = filename
        missingNote.Title = helpers.String(noteContent).TitleFromMarkdown()
        missingNote.UserID = 1
        h.db.Create(&missingNote)
        missingNotes = append(missingNotes, missingNote.ExportedFields())
      }
    }
  }

  // Build a response object, to allow more data
  var response = map[string]interface{}{ 
    "notes": missingNotes,
  }
  c.JSON(http.StatusOK, gin.H{ "data": response })
}

package controllers

import "github.com/jinzhu/gorm"

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


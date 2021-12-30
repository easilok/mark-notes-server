package models

type NoteInformation struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Filename string `json:"filename" gorm:"uniqueIndex"`
	Title    string `json:"title"`
	Favorite bool   `json:"favorite"`
	UserID   uint
	User     User
}

type Category struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
}

type NoteInformationAPI struct {
	// ID        uint    `json:"id"`
	Filename string `json:"filename"`
	Title    string `json:"title"`
	Favorite bool   `json:"favorite"`
}

func (base *NoteInformation) ExportedFields() NoteInformationAPI {
	var exported NoteInformationAPI

	exported.Title = base.Title
	exported.Filename = base.Filename
	exported.Favorite = base.Favorite

	return exported
}

package models

type NoteInformation struct {
  ID        uint    `json:"id" gorm:"primary_key"`
  Filename  string  `json:"filename" gorm:"uniqueIndex"`
  Title     string  `json:"tile"`
  Favorite  bool    `json:"favorite"`
}

type Category struct {
  ID        uint    `json:"id" gorm:"primary_key"`
  Title     string  `json:"tile"`
}

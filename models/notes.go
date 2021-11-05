package models

type User struct {
  ID        uint    `json:"id" gorm:"primary_key"`
  Email     string  `json:"email" gorm:"uniqueIndex"`
  Password  string  `json:"password"`
  Name      string  `json:"name"`
}

type NoteInformation struct {
  ID        uint    `json:"id" gorm:"primary_key"`
  Filename  string  `json:"filename" gorm:"uniqueIndex"`
  Title     string  `json:"tile"`
  Favorite  bool    `json:"favorite"`
  UserID    uint     
}

type Category struct {
  ID        uint    `json:"id" gorm:"primary_key"`
  Title     string  `json:"tile"`
}

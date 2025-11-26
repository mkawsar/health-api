package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	AuthorID uint   `json:"author_id" gorm:"not null;index"`
	Author   User   `json:"author" gorm:"foreignKey:AuthorID"`
	Title    string `json:"title" gorm:"not null"`
	Content  string `json:"content" gorm:"type:text"`
}

func NewNote(authorID uint, title string, content string) *Note {
	return &Note{
		AuthorID: authorID,
		Title:    title,
		Content:  content,
	}
}

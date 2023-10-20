package models

import "gorm.io/gorm"

// Book is a struct that represents a book in the database
type Book struct {
	gorm.Model         // gorm.Model provides the fields ID, CreatedAt, UpdatedAt, and DeletedAt
	Title       string `gorm:"size:191;not null;unique" json:"title"`
	Author      string `gorm:"size:191;not null" json:"author"`
	Publisher   string `gorm:"size:191;not null" json:"publisher"`
	Description string `gorm:"size:191;not null" json:"description"`
}

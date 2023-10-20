package models

import (
	"errors"
	"html"
	"strings"

	"github.com/zerodot618/go-huang/api/database"
	"gorm.io/gorm"
)

// ToDo is a struct contains information about a to-do item
type ToDo struct {
	gorm.Model
	Title    string `gorm:"size:100;not null;unique" json:"title"`
	Content  string `gorm:"text;not null" json:"content"`
	Author   User   `json:"author"`
	AuthorID uint   `gorm:"not null" json:"author_id"`
}

// Prepare is a method that prepares the ToDo struct for saving
// It escapes the title and content strings and sets the created and updated times
func (t *ToDo) Prepare() {
	t.Title = html.EscapeString(strings.TrimSpace(t.Title))
	t.Content = html.EscapeString(strings.TrimSpace(t.Content))
	t.Author = User{}
}

// Validate is a method that validates the ToDo struct
// It returns a map of errors if any of the fields are not valid
func (t *ToDo) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)
	if t.Title == "" {
		err = errors.New("required title")
		errorMessages["Required_title"] = err.Error()
	}
	if t.Content == "" {
		err = errors.New("required content")
		errorMessages["Required_content"] = err.Error()
	}
	if t.AuthorID < 1 {
		err = errors.New("required author")
		errorMessages["Required_author"] = err.Error()
	}
	return errorMessages
}

// SaveToDo is a method that saves a ToDo struct to the database
// It takes a pointer to a gorm.DB as an argument and returns a pointer to the saved ToDo struct
func (t *ToDo) SaveToDo(db *gorm.DB) (*ToDo, error) {
	err := database.GlobalDB.Model(&ToDo{}).Create(&t).Error
	if err != nil {
		return &ToDo{}, err
	}
	if t.Author.ID != 0 {
		err = database.GlobalDB.Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
		if err != nil {
			return &ToDo{}, err
		}
	}
	return t, nil
}

// UpdateAToDo is a method that updates a ToDo struct in the database
// It takes a pointer to a gorm.DB as an argument and returns a pointer to the updated ToDo struct
func (t *ToDo) UpdateAToDo(db *gorm.DB) (*ToDo, error) {
	err := database.GlobalDB.Model(&ToDo{}).Where("id = ?", t.ID).Updates(ToDo{Title: t.Title, Content: t.Content}).Error
	if err != nil {
		return &ToDo{}, err
	}
	if t.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", t.AuthorID).Take(&t.Author).Error
		if err != nil {
			return &ToDo{}, err
		}
	}
	return t, nil
}

// DeleteAToDo is a method that deletes a ToDo struct from the database
// It takes a pointer to a gorm.DB as an argument and returns the number of rows affected
func (t *ToDo) DeleteAToDo() (int64, error) {
	db := database.GlobalDB.Model(&ToDo{}).Where("id = ?", t.ID).Take(&ToDo{}).Delete(&ToDo{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// FindUserToDos is a method that finds all ToDo structs for a given user
// It takes a pointer to a gorm.DB and a user ID as arguments and returns a pointer to a slice of ToDo structs
func (t *ToDo) FindUserToDos(uid uint) (*[]ToDo, error) {
	var err error
	todos := []ToDo{}
	err = database.GlobalDB.Model(&ToDo{}).Where("author_id = ?", uid).Limit(100).Order("created_at desc").Find(&todos).Error
	if err != nil {
		return &[]ToDo{}, err
	}
	if len(todos) > 0 {
		for i := range todos {
			err := database.GlobalDB.Model(&User{}).Where("id = ?", todos[i].AuthorID).Take(&todos[i].Author).Error
			if err != nil {
				return &[]ToDo{}, err
			}
		}
	}
	return &todos, nil
}

package models

import (
	"errors"
	"html"
	"strings"

	"github.com/zerodot618/go-huang/api/database"
	"github.com/zerodot618/go-huang/api/security"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User defines the user in db
// User struct is used to store user information in the database
type User struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email" gorm:"unique"`
	Password string `json:"password" binding:"required"`
}

// BeforeSave is a hook that is called before a user is saved to the database
// It hashes the user's password before saving it to the database
func (user *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := security.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// Prepare is a function that is called before a user is saved to the database
// It escapes any HTML characters and trims any whitespace
func (user *User) Prepare() {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
}

// Validate is a function that is used to validate a user before saving it to the database
// It takes an action as an argument, which is used to determine which validation to perform
func (user *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error
	switch strings.ToLower(action) {
	case "update":
		if user.Email == "" {
			err = errors.New("required email")
			errorMessages["Required_email"] = err.Error()
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				err = errors.New("invalid email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	case "login":
		if user.Password == "" {
			err = errors.New("eequired password")
			errorMessages["Required_password"] = err.Error()
		}
		if user.Email == "" {
			err = errors.New("required email")
			errorMessages["Required_email"] = err.Error()
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				err = errors.New("invalid email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	default:
		if user.Name == "" {
			err = errors.New("required name")
			errorMessages["Required_name"] = err.Error()
		}
		if user.Password == "" {
			err = errors.New("required password")
			errorMessages["Required_password"] = err.Error()
		}
		if user.Password != "" && len(user.Password) < 6 {
			err = errors.New("password should be atleast 6 characters")
			errorMessages["Invalid_password"] = err.Error()
		}
		if user.Email == "" {
			err = errors.New("required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				err = errors.New("invalid email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	}
	return errorMessages
}

// SaveUser is a function that is used to save a user to the database
// It takes a pointer to a gorm.DB as an argument
func (u *User) SaveUser() (*User, error) {
	err := database.GlobalDB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// CreateUserRecord creates a user record in the database
// CreateUserRecord takes a pointer to a User struct and creates a user record in the database
// It returns an error if there is an issue creating the user record
func (user *User) CreateUserRecord() error {
	result := database.GlobalDB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// HashPassword encrypts user password
// HashPassword takes a string as a parameter and encrypts it using bcrypt
// It returns an error if there is an issue encrypting the password
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword checks user password
// CheckPassword takes a string as a parameter and compares it to the user's encrypted password
// It returns an error if there is an issue comparing the passwords
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

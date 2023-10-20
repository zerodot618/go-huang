package models

import (
	"math/rand"
	"time"

	"github.com/zerodot618/go-huang/api/database"
	"gorm.io/gorm"
)

// URL is a struct that stores the information for a URL
type URL struct {
	gorm.Model
	LongURL      string     `json:"long_url" gorm:"unique"`
	ShortURL     string     `json:"short_url" gorm:"unique"`
	AccessCount  uint       `json:"access_count"`
	LastAccessed *time.Time `json:"last_accessed"`
	AccessPlace  string     `json:"access_place"`
}

// GenerateShortURL is a method used to generate a random short URL
// It takes no parameters and returns nothing
func (u *URL) GenerateShortURL() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	u.ShortURL = string(b)
}

// CreateURL is a method used to create a new URL in the database
// It takes in a pointer to a URL struct as a parameter and returns nothing
func CreateURL(url *URL) {
	database.GlobalDB.Create(&url)
}

// GetURLByShortURL is a method used to get a URL from the database by its short URL
// It takes a string as a parameter and returns a URL struct and an error
func GetURLByShortURL(shortURL string) (URL, error) {
	var url URL
	if err := database.GlobalDB.Where("short_url = ?", shortURL).First(&url).Error; err != nil {
		return url, err
	}
	return url, nil
}

// UpdateURL is a method used to update a URL in the database
func UpdateURL(url *URL) error {
	result := database.GlobalDB.Save(url)
	return result.Error
}

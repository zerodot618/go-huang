package database

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GlobalDB is a global db object that will be used across different packages
var GlobalDB *gorm.DB

// InitDatabase creates a mysql db connecntion and stores it in the GlobalDB variable
// It reads the environment variables from the .env file and uses them to create the connection
// It returns an error if the connection fails
func InitDatabase() (err error) {
	// Read the environment variables from the .env file
	config, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error reading .env file")
	}
	// Get database connection details from environment variables
	dbHost := config["DB_HOST"]
	dbPort := config["DB_PORT"]
	dbUser := config["DB_USER"]
	dbPass := config["DB_PASS"]
	dbName := config["DB_NAME"]

	// Creatre the data source name (DSN) using the environment variables
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)
	// Create the connection and store it in the GlobalDB variable
	GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return
}

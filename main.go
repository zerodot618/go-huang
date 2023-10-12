package main

import (
	"log"

	"github.com/zerodot618/go-huang/database"
	"github.com/zerodot618/go-huang/models"
	"github.com/zerodot618/go-huang/routes"
)

// main is the entry point of the program.
// It initializes the database, sets up the router and starts the server.
func main() {
	// Initialize the database
	err := database.InitDatabase()
	if err != nil {
		// Log the error and exit
		log.Fatalln("could not create database")
	}
	// Automigrate the User model
	// AutoMigrate() automatically migrates our schema, to keep our schema upto date.
	database.GlobalDB.AutoMigrate(&models.User{})
	database.GlobalDB.AutoMigrate(&models.Book{})
	// Set up the router
	r := routes.SetupRouter()
	// Start the server
	r.Run(":8088")
}

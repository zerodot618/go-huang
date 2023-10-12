package main

import (
	"log"

	"github.com/zerodot618/go-huang/database"
	"github.com/zerodot618/go-huang/models"
	"github.com/zerodot618/go-huang/routes"

	_ "github.com/zerodot618/go-huang/docs"
)

// main is the entry point of the program.
// It initializes the database, sets up the router and starts the server.

// @title Swagger JWT API
// @version 1.0
// @description Create  Go REST API with JWT Authentication in Gin Framework
// @contact.name API Support
// @termsOfService demo.com
// @contact.url http://demo.com/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8088
// @BasePath /api
// @Schemes http https
// @query.collection.format multi
// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
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
	database.GlobalDB.AutoMigrate(&models.URL{})
	database.GlobalDB.AutoMigrate(&models.File{})
	// Set up the router
	r := routes.SetupRouter()
	// Start the server
	r.Run(":8088")
}

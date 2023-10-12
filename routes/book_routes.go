package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zerodot618/go-huang/controllers"
)

func setupBookRoutes(router *gin.RouterGroup) {
	/*
	   This function sets up the routes for the book-related API endpoints.
	   It takes in a pointer to a gin.Engine instance and a pointer to a gorm.DB instance.
	   It creates a new instance of the BookController and sets up the routes for the API endpoints.
	*/
	// Create a new instance of the BookController with the provided gorm.DB instance
	var bookController controllers.BookController

	// Create a new group of routes for the book-related API endpoints
	bookRoutes := router.Group("/books")
	{
		// Set up the route for creating a new book
		bookRoutes.POST("", bookController.CreateBook)

		// Set up the route for searching books
		bookRoutes.GET("/search", bookController.SearchBooks)

		// Set up the route for getting a list of all books
		bookRoutes.GET("", bookController.GetBooksList)

		// Set up the route for getting a book by its ID
		bookRoutes.GET("/:id", bookController.GetBookById)

		// Set up the route for updating a book by its ID
		bookRoutes.PUT("/:id", bookController.UpdateBook)

		// Set up the route for deleting a book by its ID
		bookRoutes.DELETE("/:id", bookController.DeleteBook)
	}
}

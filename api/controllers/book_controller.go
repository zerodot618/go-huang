package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zerodot618/go-huang/api/database"
	"github.com/zerodot618/go-huang/api/models"
)

// BookController is a struct that represents a controller for book-related operations
type BookController struct{}

// CreateBook is a function that creates a new book record in the database
func (ctrl *BookController) CreateBook(c *gin.Context) {
	// Bind the request body to a Book struct
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create a new book record in the database
	if err := database.GlobalDB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the created book to the client
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// SearchBooks is a function that searches the database for books that match a given query
func (ctrl *BookController) SearchBooks(c *gin.Context) {
	// Get the search query from the request parameters
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}
	// Query the database for books that match the search query
	var books []models.Book
	if err := database.GlobalDB.Where("title LIKE ? OR author LIKE ? OR publisher LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the search results to the client
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// DeleteBook is a function that deletes a book record from the database
func (ctrl *BookController) DeleteBook(c *gin.Context) {
	// Get the book ID from the request parameters
	bookID := c.Param("id")
	// Check if the book with the given ID exists
	var book models.Book
	if err := database.GlobalDB.Where("id = ?", bookID).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	// Delete the book record from the database
	if err := database.GlobalDB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return a success message to the client
	c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
}

// UpdateBook is a function that updates a book record in the database
func (ctrl *BookController) UpdateBook(c *gin.Context) {
	// Get the book ID from the request parameters
	id := c.Param("id")
	// Check if the book with the given ID exists
	var book models.Book
	if err := database.GlobalDB.Where("id = ?", id).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	// Bind the request body to a struct
	var updateData struct {
		Title       *string `json:"title"`
		Author      *string `json:"author"`
		Publisher   *string `json:"publisher"`
		Description *string `json:"description"`
	}
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update the book fields with the provided values
	if updateData.Title != nil {
		book.Title = *updateData.Title
	}
	if updateData.Author != nil {
		book.Author = *updateData.Author
	}
	if updateData.Publisher != nil {
		book.Publisher = *updateData.Publisher
	}
	if updateData.Description != nil {
		book.Description = *updateData.Description
	}
	// Save the updated book to the database
	if err := database.GlobalDB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the updated book to the client
	c.JSON(http.StatusOK, gin.H{"data": book})
}

// GetBooksList is a function that retrieves a list of all book records from the database
func (ctrl *BookController) GetBooksList(c *gin.Context) {
	// Query the database for all book records
	var books []models.Book
	if err := database.GlobalDB.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return the list of books to the client
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// GetBookById is a function that retrieves a single book record from the database by ID
func (ctrl *BookController) GetBookById(c *gin.Context) {
	// Get the book ID from the request parameters
	bookID := c.Param("id")
	// Check if the book with the given ID exists
	var book models.Book
	if err := database.GlobalDB.Where("id = ?", bookID).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	// Return the book to the client
	c.JSON(http.StatusOK, gin.H{"data": book})
}

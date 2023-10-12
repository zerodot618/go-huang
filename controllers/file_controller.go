package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zerodot618/go-huang/database"
	"github.com/zerodot618/go-huang/models"
)

// FileController is a struct that represents a controller for file-related operations“
type FileController struct{}

// UploadFile is a function that handles the upload of a single file
func (f *FileController) UploadFile(c *gin.Context) {
	/*
	 UploadFile function handles the upload of a single file.
	 It gets the file from the form data, saves it to the defined path,
	 generates a unique identifier for the file, saves the file metadata tp the database,
	 and returns a success message and the file metadata.
	*/
	// Get the file from the form data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	// Define the path where the file will be saved
	filePath := filepath.Join("uploads", file.Filename)
	// Save the file to the defined path
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})

		return
	}
	// Generate a unique identifier for the file
	uuid := uuid.New().String()
	// Save file metadata to database
	fileMetadata := models.File{
		Filename: file.Filename,
		UUID:     uuid,
	}
	if err := database.GlobalDB.Create(&fileMetadata).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})

		return
	}
	// Return success message and file metadata
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    fileMetadata,
	})
}

// UploadFiles is a function that handles the upload of multiple files
func (f *FileController) UploadFiles(c *gin.Context) {
	/*
	  UploadFiles function handles the upload of multiple files.
	  It gets the files from the form data, saves each file to the defined path,
	  generates a unique identifier for each file, saves the file metadata to the database,
	  and returns a success message and the file metadata.
	*/
	// Get the files from the form data
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	files := form.File["files"]
	var fileModels []models.File
	// Save each file to the defined path and generate a unique identifier for each file
	for _, file := range files {
		filePath := filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})

			return
		}
		fileModels = append(fileModels, models.File{
			UUID:     uuid.New().String(),
			Filename: file.Filename,
		})
	}
	// Save file metadata to database
	err = database.GlobalDB.Create(&fileModels).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})

		return
	}
	// Return a success message and the file metadata
	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
		"files":   fileModels,
	})
}

// GetFile is a function that retrieves a file from the server
func (f *FileController) GetFile(c *gin.Context) {
	/*
	 GetFile function retrieves a file from the server.
	 It gets the unique identifier of the file to be retrieved,
	 retrieves the file metadata from the database,
	 defines the path of the file to be retrieved,
	 opens the file, reads the first 512 bytes of the file to determine its content type,
	 gets the file info, sets the headers for the file transfer, and returns the file.
	*/
	// Get the unique identifier of the file to be retrieved
	uuid := c.Param("uuid")
	var file models.File
	// Retrieve the file metadata from the database
	err := database.GlobalDB.Where("uuid = ?", uuid).First(&file).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	// Define the path of the file tobe retrieved
	filePath := filepath.Join("uploads", file.Filename)
	// Open the file
	fileData, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileData.Close()
	// Read the first 512 bytes of the file to determine its content type
	fileHeader := make([]byte, 512)
	_, err = fileData.Read(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
	}
	fileContentType := http.DetectContentType(fileHeader)
	// Get the file info
	fileInfo, err := fileData.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}
	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	c.Header("Content-Type", fileContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.File(filePath)
}

// DeleteFile is a function that deletes a file from the server and its metadata from the database
func (f *FileController) DeleteFile(c *gin.Context) {
	/*
	 DeleteFile function deletes a file from the server and its metadata from the database.
	 It gets the unique identifier of the file to be deleted,
	 retrieves the file metadata from the database,
	 defines the path of the file to be deleted,
	 deletes the file from the server,
	 deletes the file metadata from the database,
	 and returns a success message.
	*/
	// Get the unique identifier of the file to be deleted
	uuid := c.Param("uuid")
	var file models.File
	// Retrieve the file metadata from the database
	err := database.GlobalDB.Where("uuid = ?", uuid).First(&file).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	// Define the path of the file to be deleted
	filePath := filepath.Join("uploads", file.Filename)
	// Delete the file from the server
	err = os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}
	// Delete the file metadata from the database
	err = database.GlobalDB.Delete(&file).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from database"})
		return
	}
	// Return a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "File " + file.Filename + " deleted successfully",
	})
}

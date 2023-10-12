package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zerodot618/go-huang/controllers"
)

func setupFileRoutes(router *gin.RouterGroup) {
	var fileController controllers.FileController

	fileRoutes := router.Group("/files")
	{
		fileRoutes.POST("/file", fileController.UploadFile)
		fileRoutes.POST("/files", fileController.UploadFiles)
		fileRoutes.GET("/file/:uuid", fileController.GetFile)
		fileRoutes.DELETE("/file/:uuid", fileController.DeleteFile)
	}
}

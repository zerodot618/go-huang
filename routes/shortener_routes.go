package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zerodot618/go-huang/controllers"
)

func setupShortenerRoutes(router *gin.RouterGroup) {
	var shortenerController controllers.ShortenerController

	shortenerRoutes := router.Group("/shortener")
	{
		shortenerRoutes.POST("", shortenerController.CreateShortURL)
		shortenerRoutes.GET("/:short_url", shortenerController.RedirectShortURL)
		shortenerRoutes.GET("/:short_url/stats", shortenerController.GetURLStatistics)
	}
}

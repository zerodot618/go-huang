package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zerodot618/go-huang/controllers"
	"github.com/zerodot618/go-huang/middlewares"
)

func setupUserRoutes(router *gin.RouterGroup) {
	var userController controllers.UserController

	// Create a new group for the public routes
	public := router.Group("/public")
	{
		// Add the login route
		public.POST("/login", userController.Login)
		// Add the signup route
		public.POST("/signup", userController.Signup)
	}

	// Add the signup route
	protected := router.Group("/protected").Use(middlewares.Authz())
	{
		// Add the profile route
		protected.GET("/profile", userController.Profile)
	}
}

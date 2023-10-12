package routes

import (
	"github.com/gin-gonic/gin"
)

// setupRouter sets up the router and adds the routes.
func SetupRouter() *gin.Engine {
	// Create a new router
	r := gin.Default()
	// Add a welcome route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome To This Website")
	})
	// Create a new group for the API
	api := r.Group("/api")
	{
		// Add the routes for the user
		setupUserRoutes(api)
		setupBookRoutes(api)
	}
	// Return the router
	return r
}

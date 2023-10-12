package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zerodot618/go-huang/auth"
)

// Authz is a middleware that validates token and authorizes users
// It takes a gin.Context as an argument and returns a gin.HandlerFunc
// This function is responsible for validating the token sent by the client in the Authorization header
// and authorizing the user if the token is valid
func Authz() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header from the request
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			// If the Authorization header is not present, return a 403 status code
			c.JSON(http.StatusForbidden, "No Authorization header provided")
			c.Abort()
			return
		}
		// Split the Authorization header to get the token
		extractedToken := strings.Split(clientToken, "Bearer ")
		if len(extractedToken) == 2 {
			// Trime the token
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			// If the token is not in the correct format, return a 400 status code
			c.JSON(http.StatusBadRequest, "Incorrect Format of Authorization Token")
			c.Abort()
			return
		}
		// Create a JwtWrapper with the secret key and issuer
		jwtWrapper := auth.JwtWrapper{
			SecretKey: "verysecretkey",
			Issuer:    "AuthService",
		}
		// Validate the token
		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			// If token is not valid, return a 401 status code
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		// Set the claims in the context
		c.Set("email", claims.Email)
		// Continue to the next handler
		c.Next()
	}
}

package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zerodot618/go-huang/api/auth"
	"github.com/zerodot618/go-huang/api/database"
	"github.com/zerodot618/go-huang/api/models"
	"gorm.io/gorm"
)

// UserController is a struct that represents a controller for user-related operations
type UserController struct{}

// LoginPayload login body
// LoginPayload is a struct that contains the fields for a user's login credentials
type LoginPayload struct {
	Email    string `json:"email" bingding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse token response
// LoginResponse is a struct that contains the fields for a user's login response
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// swagger input
type UserDetails struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Signup is a function that handles user signup
// It takes in a gin context as an argument and binds the user data from the request body to a user struct
// It then hashes the user's password and creates a user record in the database
// If successful, it returns a 200 status code with a success message
// If unsuccessful, it returns a 400 or 500 status code with an error message

// @Summary Signup User
// @Description Signin
// @ID SignupUser
// @Tags User
// @Param EnterDetails body UserDetails true "Signin"
// @Accept json
// @Success 200 {object} string "Success"
// @Failure 400 {string} string "Error"
// @Router /public/signup [POST]
func (ctrl UserController) Signup(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Inputs"})
		c.Abort()
		return
	}
	err = user.HashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Hashing Password"})
		c.Abort()
		return
	}
	err = user.CreateUserRecord()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error Creating User"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Sucessfully Register"})
}

// Login is a function that handles user login
// It takes in a gin context as an argument and binds the user data from the request body to a LoginPayload struct
// It then checks if the user exists in the database and if the password is correct
// If successful, it generates a token and a refresh token and returns a 200 status code with the token and refresh token
// If unsuccessful, it returns a 401 or 500 status code with an error message

// @Summary Login User
// @Description  Login
// @Tags User
// @ID LoginUser
// @Param EnterDetails body LoginPayload true "Login"
// @Accept json
// @Success 200  {object}  string  "Success"
// @Failure 400  {string}  string  "Error"
// @Router /public/login [POST]
func (ctrl UserController) Login(c *gin.Context) {
	var payload LoginPayload
	var user models.User
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Inputs"})
		c.Abort()
		return
	}
	result := database.GlobalDB.Where("email = ?", payload.Email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid User Credentials"})
		c.Abort()
		return
	}
	err = user.CheckPassword(payload.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid User Credentials"})
		c.Abort()
		return
	}
	signedToken, err := auth.CreateToken(user)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		c.Abort()
		return
	}
	tokenResponse := LoginResponse{
		Token: signedToken,
	}
	c.JSON(200, tokenResponse)
}

// Profile is a controller function that retrieves the user profile from the database
// based on the email provided in the authorization middleware.
// It returns a 404 status code if the user is not found,
// and a 500 status code if an error occurs while retrieving the user profile.

// @Summary Get User By Token
// @ID GetUserByToken
// @Produce json
// @Tags User
// Param [param_name] [param_type] [data_type] [required/mandatory] [description]
// @Param Authorization header string true "Authorization header using the Bearer scheme"
// @Success 200 {object} string "Success"
// @Failure 400 {string} string "Error"
// @Router /protected/profile [GET]
func (ctrl UserController) Profile(c *gin.Context) {
	// Initialize a user model
	var user models.User
	// Get the eamil from the authorization middleware
	email, _ := c.Get("email")
	// Query the database for the user
	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	// If the user is not found, return a 404 status code
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User Not Found"})
		c.Abort()
		return
	}
	// If an error occurs while retriving the user profile, return a 500 status code
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Could Not Get User Profile"})
		c.Abort()
		return
	}
	// Set the user's password to an empty string
	user.Password = ""
	// Return the user profile with a 200 status code
	c.JSON(http.StatusOK, user)
}

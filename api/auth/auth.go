package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/zerodot618/go-huang/api/models"
)

// JwtClaim adds email as a claim to the token
// JwtClaim is a struct that holds the Email of the user, as well as the StandardClaims
type JwtClaim struct {
	ID         uint
	Email      string
	Authorized bool
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token
// GenerateToken takes an email as an argument and returns a signed JWT token and an error
func CreateToken(user models.User) (signedToken string, err error) {
	claims := &JwtClaim{
		ID:         user.ID,
		Email:      user.Email,
		Authorized: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Minute * time.Duration(120))),
			Issuer:    "AuthService",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return
	}
	return
}

// RefreshToken generates a refresh jwt token
// RefreshToken takes an email as an argument and returns a signed JWT token and an error
func RefreshToken(user models.User) (signedtoken string, err error) {
	claims := &JwtClaim{
		ID:         user.ID,
		Email:      user.Email,
		Authorized: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Minute * time.Duration(120))),
			Issuer:    "AuthService",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedtoken, err = token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return
	}
	return
}

// ValidateToken validates the JWT token
// ValidateToken takes a signed JWT token as an argument and returns the JwtClaim and an error
func ValidateToken(r *http.Request) error {
	tokenString := Extractoken(r)
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("API_SECRET")), nil
		},
	)
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(*JwtClaim); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

// ExtractToken extracts the token from the http request
// r: the http request
// returns: the token string
func Extractoken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenID extracts the token id from the http request
// r: the http request
// returns: the token id and an error if any
func ExtractTokenID(r *http.Request) (uint, error) {
	tokenString := Extractoken(r)
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("API_SECRET")), nil
		},
	)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*JwtClaim)
	if ok && token.Valid {
		return claims.ID, nil
	}
	return 0, err
}

// Pretty display the claims licely in the terminal
// data: the data to be displayed
func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", "")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}

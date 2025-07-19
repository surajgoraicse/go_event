package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/surajgoraicse/go_event/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
type loginResponse struct {
	Token string `json:"token"`
}

func (app *application) login(c *gin.Context) {
	req := &loginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	existingUser, err := app.models.Users.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err, "message": "user not found, please signup"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID" : existingUser.ID,
		"expr" : time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" :"error generating token"})
		return 
	}
	c.JSON(http.StatusOK, loginResponse{Token: tokenString})



}

func (app *application) registerUser(c *gin.Context) {
	register := new(registerRequest)

	if err := c.ShouldBindJSON(register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("internal server error")
		c.JSON(http.StatusInternalServerError, gin.H{"bcrypt error": err})
		return
	}
	register.Password = string(hashedPassword)
	user := &database.User{
		Email:    register.Email,
		Name:     register.Name,
		Password: register.Password,
	}

	if err := app.models.Users.Insert(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, user) // will not contain the password in the response

}

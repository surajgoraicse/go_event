package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/surajgoraicse/go_event/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
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

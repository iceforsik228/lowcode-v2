package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"lowcode-v2/initializers"
	"lowcode-v2/models"
	"lowcode-v2/utils"
	"net/http"
	"os"
	"time"
)

type registerInput struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func Register(c *gin.Context) {
	// get data from request
	var input registerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// add user to db
	user := models.User{Email: input.Email, Password: string(hashedPassword)}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// send a response
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func Login(c *gin.Context) {
	// get data from request
	var input registerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// get user by data (check id too)
	var user models.User
	result := initializers.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil || user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find user",
		})
		return
	}

	// check password & hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to compare passwords (password is wrong!)",
		})
		return
	}

	// generate jwt
	token := utils.GenerateToken(user.ID, time.Hour)

	// Sign and get the complete encoded token as a string using the secret
	hmacSecret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	// return jwt token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Logout(c *gin.Context) {
	token := utils.GenerateToken(0, -time.Hour)

	hmacSecret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func Validate(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		panic("yamaha > kawasaki")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "I'm logged in!",
		"user":    user,
	})
}

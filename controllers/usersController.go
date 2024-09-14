package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/LoTfI01101011/go_blog/initial"
	"github.com/LoTfI01101011/go_blog/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("Secret")))
}

func SignUpUser(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
	}
	//hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash the password",
		})
	}
	//creating the new user
	id, _ := uuid.NewV7()
	user := models.User{ID: id, Email: body.Email, Password: string(hash)}
	initial.DB.Create(&user)
	//generating the jwt token
	token, err := GenerateToken(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create the token",
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

func LoginUser(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//geting the user from the db
	var user models.User
	if err := initial.DB.Where("email = ?", body.Email).Find(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "The provided credentials are incorrect",
		})
		return
	}
	//comparing the two paaswords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "The provided password is incorrect",
		})
		return
	}
	//generating the jwt token
	token, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create the token",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

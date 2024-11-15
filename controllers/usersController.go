package controllers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LoTfI01101011/Library/initial"
	"github.com/LoTfI01101011/Library/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("Secret")))
}
func addTokenToBlacklist(token string, rdb *redis.Client, ctx context.Context) error {
	_, err := rdb.Set(ctx, token, "blacklisted", time.Hour*24).Result()
	return err
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
func Logout(c *gin.Context) {
	//geting the token from the header
	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	//adding the token to the blacklist
	err := addTokenToBlacklist(token, initial.Rdb, initial.Ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.Status(200)

}

func BeginOAuthHundler(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}
func CallbackAuthHundler(c *gin.Context) {
	//get the user info
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
	}
	//generate the token
	// userUUID, _ := uuid.Parse(user.UserID)
	// token, err := GenerateToken(userUUID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err,
	// 	})
	// }
	//insert the user in the db
	id, _ := uuid.NewV7()
	User := models.User{ID: id, Email: user.Email}
	initial.DB.Create(&User)
	//this a test of how do i pass the token to the next page
	// http.Redirect(c.Writer, c.Request, "https://lotfi-portfolioo.netlify.app/?token="+token, http.StatusFound)

}

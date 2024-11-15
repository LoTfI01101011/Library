package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LoTfI01101011/Library/initial"
	"github.com/LoTfI01101011/Library/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddelware(c *gin.Context) {
	//get jwt token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	//check in blackedlist
	isBlacklisted, err := initial.Rdb.Get(initial.Ctx, tokenString).Result()
	if err == nil && isBlacklisted == "blacklisted" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	//decode/validate it
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		return []byte(os.Getenv("Secret")), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//check the expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//get the user from the db with token sub
		var user models.User
		err := initial.DB.Where("id = ?", claims["sub"]).Find(&user).Error
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//set the user
		c.Set("user", user)
		//continue
		c.Next()
	}

}

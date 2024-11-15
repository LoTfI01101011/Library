package main

import (
	"os"

	"github.com/LoTfI01101011/Library/controllers"
	"github.com/LoTfI01101011/Library/initial"
	"github.com/LoTfI01101011/Library/middleware"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func init() {
	initial.ConnectToDb()
	initial.InitRedis()
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:8000/api/google/auth/callback?provider=google"),
	)
}

func main() {

	r := gin.Default()
	//user
	r.POST("/api/auth/register", controllers.SignUpUser)
	r.POST("/api/auth/login", controllers.LoginUser)
	r.POST("/api/auth/logout", middleware.AuthMiddelware, controllers.Logout)
	r.GET("/api/auth/:provider", controllers.BeginOAuthHundler)
	r.GET("/api/:provider/auth/callback", controllers.CallbackAuthHundler)
	//book
	r.POST("/api/book", middleware.AuthMiddelware, controllers.CreateBook)
	r.GET("/api/book", middleware.AuthMiddelware, controllers.GetBooks)
	r.GET("/api/book/:id", middleware.AuthMiddelware, controllers.GetBookById)
	r.PATCH("/api/book/:id", middleware.AuthMiddelware, controllers.UpdateBook)
	r.DELETE("/api/book/:id", middleware.AuthMiddelware, controllers.DeleteBook)
	r.Run()
}

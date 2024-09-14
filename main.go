package main

import (
	"github.com/LoTfI01101011/go_blog/controllers"
	"github.com/LoTfI01101011/go_blog/initial"
	"github.com/LoTfI01101011/go_blog/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initial.ConnectToDb()
}

func main() {
	r := gin.Default()
	//user
	r.POST("/api/register", controllers.SignUpUser)
	r.POST("/api/login", controllers.LoginUser)
	//book
	r.POST("/api/book", middleware.AuthMiddelware, controllers.CreateBook)
	r.GET("/api/book", middleware.AuthMiddelware, controllers.GetBooks)
	r.GET("/api/book/:id", middleware.AuthMiddelware, controllers.GetBookById)
	r.PATCH("/api/book/:id", middleware.AuthMiddelware, controllers.UpdateBook)
	r.DELETE("/api/book/:id", middleware.AuthMiddelware, controllers.DeleteBook)
	r.Run()
}

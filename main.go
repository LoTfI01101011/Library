package main

import (
	"github.com/LoTfI01101011/go_blog/controllers"
	"github.com/LoTfI01101011/go_blog/initial"
	"github.com/gin-gonic/gin"
)

func init() {
	initial.ConnectToDb()
}

func main() {
	r := gin.Default()
	r.POST("/book", controllers.CreateBook)
	r.GET("/book", controllers.GetBooks)
	r.GET("/book/:id", controllers.GetBookById)
	r.PATCH("/book/:id", controllers.UpdateBook)
	r.DELETE("/book/:id", controllers.DeleteBook)
	r.Run()
}

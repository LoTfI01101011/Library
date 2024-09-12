package main

import (
	"github.com/LoTfI01101011/go_blog/initial"
	"github.com/gin-gonic/gin"
)

func init() {
	initial.ConnectToDb()
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

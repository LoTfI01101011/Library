package controllers

import (
	"github.com/LoTfI01101011/go_blog/initial"
	"github.com/LoTfI01101011/go_blog/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateBook(c *gin.Context) {

	var body struct {
		title       string
		author      string
		pages       int32
		description string
	}
	//get data from the request body
	c.Bind(&body)

	// create a book
	id, _ := uuid.NewV7()
	book := models.Book{Id: id, Title: body.title, Author: body.author, Pages: int(body.pages), Description: body.description}
	result := initial.DB.Create(&book)
	if result.Error != nil {
		c.Status(400)
		return
	}
	//returning the instance
	c.JSON(200, gin.H{
		"book": book,
	})

}
func GetBooks(c *gin.Context) {
	var books []models.Book

	initial.DB.Find(&books)

	c.JSON(200, gin.H{
		"books": books,
	})
}

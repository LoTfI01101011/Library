package controllers

import (
	"github.com/LoTfI01101011/Library/initial"
	"github.com/LoTfI01101011/Library/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateBook(c *gin.Context) {

	var body struct {
		Title       string `json:"Title"`
		Author      string `json:"Author"`
		Pages       int32  `json:"Pages"`
		Description string `json:"Description"`
	}
	//get data from the request body
	c.Bind(&body)
	//get the user from the token
	user, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(404)
	}
	User := user.(models.User)
	// create a book
	id, _ := uuid.NewV7()
	book := models.Book{ID: id, UserID: User.ID, Title: body.Title, Author: body.Author, Pages: int(body.Pages), Description: body.Description}
	result := initial.DB.Create(&book)
	if result.Error != nil {
		c.Status(400)
		return
	}
	//reload the book with the associate user
	initial.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email")
	}).First(&book, book.ID)
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

func GetBookById(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	initial.DB.Where("id = ?", id).Find(&book)

	c.JSON(200, gin.H{
		"book": book,
	})
}
func UpdateBook(c *gin.Context) {
	//get the id from the url
	id := c.Param("id")

	//get the data from the body
	var body struct {
		Title       string `json:"Title"`
		Author      string `json:"Author"`
		Pages       int32  `json:"Pages"`
		Description string `json:"Description"`
	}

	c.Bind(&body)

	//find  the book
	var book models.Book
	err := initial.DB.Where("id = ?", id).Find(&book).Error
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}
	//update the book
	//validation
	updateData := map[string]interface{}{}
	if body.Title != "" {
		updateData["Title"] = body.Title
	}
	if body.Author != "" {
		updateData["Author"] = body.Author
	}
	if body.Pages != 0 {
		updateData["Pages"] = body.Pages
	}
	if body.Description != "" {
		updateData["Description"] = body.Description
	}
	initial.DB.Model(&book).Updates(updateData)
	// return the insctance of the book
	c.JSON(200, gin.H{
		"book": book,
	})
}
func DeleteBook(c *gin.Context) {
	//get the id from the prams
	id := c.Param("id")
	//find the book
	var book models.Book
	err := initial.DB.Where("id = ?", id).Find(&book).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}

	//delete it
	initial.DB.Delete(book)
	//return message
	c.JSON(200, gin.H{
		"response": "the book was deleted succesfuly",
	})
}

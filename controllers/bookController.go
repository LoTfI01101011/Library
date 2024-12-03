package controllers

import (
	"net/http"

	"github.com/LoTfI01101011/Library/initial"
	"github.com/LoTfI01101011/Library/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateBook(c *gin.Context) {
	var body struct {
		Title       string `json:"title"`
		Author      string `json:"author"`
		Pages       int32  `json:"pages"`
		Description string `json:"description"`
	}

	// Bind request body and check for errors
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user from the token
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Safely assert user type
	User, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
		return
	}

	// Generate a new UUID
	id, err := uuid.NewV7()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate unique ID"})
		return
	}

	// Create the book instance
	book := models.Book{
		ID:          id,
		UserID:      User.ID,
		Title:       body.Title,
		Author:      body.Author,
		Pages:       int(body.Pages),
		Description: body.Description,
	}

	// Save to database
	result := initial.DB.Create(&book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	// Preload associated user and fetch the book
	initial.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email")
	}).First(&book, "id = ?", book.ID)

	// Return the created book
	c.JSON(http.StatusOK, gin.H{"book": book})
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

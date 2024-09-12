package main

import (
	"github.com/LoTfI01101011/go_blog/initial"
	"github.com/LoTfI01101011/go_blog/models"
)

func init() {
	initial.ConnectToDb()
}

func main() {
	initial.DB.AutoMigrate(&models.Book{})
}

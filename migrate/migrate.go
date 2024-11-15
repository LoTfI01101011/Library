package main

import (
	"os"

	"github.com/LoTfI01101011/Library/initial"
	"github.com/LoTfI01101011/Library/models"
)

func init() {
	initial.ConnectToDb()
}

func main() {
	if os.Args[1] == "migrate" {
		initial.DB.AutoMigrate(&models.Book{})
	}
	if os.Args[1] == "fresh" {
		initial.DB.Migrator().DropTable(&models.Book{})
		initial.DB.AutoMigrate(&models.Book{})
	}
}

package models

import "github.com/google/uuid"

type Book struct {
	Id          uuid.UUID
	Title       string
	Author      string
	Pages       int
	Description string
}

package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	User        User
	Title       string
	Author      string
	Pages       int
	Description string
	ActivatedAt sql.NullTime // Uses sql.NullTime for nullable time fields
	CreatedAt   time.Time    // Automatically managed by GORM for creation time
	UpdatedAt   time.Time
}

package user

import (
	"time"

	"github.com/google/uuid"
)

// Declared schema
/*
	{
	"id": "d2a7924e-765f-4949-bc4c-219c956d0f8b",
	"first_name": "Alice",
	"last_name": "Bob",
	"nickname": "AB123",
	"password": "supersecurepassword",
	"email": "alice@bob.com",
	"country": "UK",
	"created_at": "2019-10-12T07:20:50.52Z",
	"updated_at": "2019-10-12T07:20:50.52Z"
	}
*/

type Model struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

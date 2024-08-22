package model

import "time"

// User is the object with user information
type User struct {
	ID        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserUpdate is the object with updating parameters
type UserUpdate struct {
	ID    int64
	Name  string
	Email string
}

// UserCreate is the object with creating parameters
type UserCreate struct {
	Name     string
	Email    string
	Password string
}

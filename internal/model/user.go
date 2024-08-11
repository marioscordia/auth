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

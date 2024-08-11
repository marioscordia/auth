package model

import "time"

// UserDB is the object that is used at the repository layer
type UserDB struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

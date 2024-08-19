package repo

import (
	"context"

	"github.com/marioscordia/auth/internal/model"
)

// Repository is an interface through which Service layer communicates with database
type Repository interface {
	CreateUser(ctx context.Context, user *model.UserCreate) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) error
	DeleteUser(ctx context.Context, id int64) error
}

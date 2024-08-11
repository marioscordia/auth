package service

import (
	"context"

	"github.com/marioscordia/auth/internal/model"
)

// Service is an interface through which Handler layer communicates with business layer
type Service interface {
	CreateUser(ctx context.Context, user *model.User, password string) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int64) error
}

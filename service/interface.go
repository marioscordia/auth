package service

import (
	"context"

	"github.com/marioscordia/auth/pkg/auth_v1"
)

// Service is ...
type Service interface {
	CreateUser(ctx context.Context, u *auth_v1.CreateRequest) (int64, error)
	GetUser(ctx context.Context, id int64) (*auth_v1.GetResponse, error)
	UpdateUser(ctx context.Context, u *auth_v1.UpdateRequest) error
	DeleteUser(ctx context.Context, id int64) error
}

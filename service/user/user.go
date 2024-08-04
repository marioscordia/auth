package user

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/marioscordia/auth/pkg/auth_v1"
	repo "github.com/marioscordia/auth/repository"
	"github.com/marioscordia/auth/service"
)

// New is ...
func New(repo repo.Repository) service.Service {
	return &serve{
		repo: repo,
	}
}

type serve struct {
	repo repo.Repository
}

func (u *serve) GetUser(ctx context.Context, id int64) (*auth_v1.GetResponse, error) {
	return u.repo.GetUser(ctx, id)
}

func (u *serve) CreateUser(ctx context.Context, req *auth_v1.CreateRequest) (int64, error) {
	hashCode := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password)))

	req.Password = hashCode

	return u.repo.CreateUser(ctx, req)
}

func (u *serve) UpdateUser(ctx context.Context, req *auth_v1.UpdateRequest) error {
	return u.repo.UpdateUser(ctx, req)
}

func (u *serve) DeleteUser(ctx context.Context, id int64) error {
	return u.repo.DeleteUser(ctx, id)
}

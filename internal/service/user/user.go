package user

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/marioscordia/auth/internal/model"
	repo "github.com/marioscordia/auth/internal/repository"
	"github.com/marioscordia/auth/internal/service"
)

// New is the function that returns Service object
func New(repo repo.Repository) service.Service {
	return &serve{
		repo: repo,
	}
}

type serve struct {
	repo repo.Repository
}

func (u *serve) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return u.repo.GetUser(ctx, id)
}

func (u *serve) CreateUser(ctx context.Context, user *model.User, password string) (int64, error) {
	hashCode := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	return u.repo.CreateUser(ctx, user, hashCode)
}

func (u *serve) UpdateUser(ctx context.Context, user *model.User) error {
	return u.repo.UpdateUser(ctx, user)
}

func (u *serve) DeleteUser(ctx context.Context, id int64) error {
	return u.repo.DeleteUser(ctx, id)
}

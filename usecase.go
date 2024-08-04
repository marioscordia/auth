package user

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/marioscordia/auth/pkg/auth_v1"
)

// New is ...
func New(repo Repository) UseCase {
	return &useCase{
		repo: repo,
	}
}

type useCase struct {
	repo Repository
}

func (u *useCase) GetUser(ctx context.Context, id int64) (*auth_v1.GetResponse, error) {
	return u.repo.GetUser(ctx, id)
}

func (u *useCase) CreateUser(ctx context.Context, req *auth_v1.CreateRequest) (int64, error) {
	hashCode := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password)))

	req.Password = hashCode

	return u.repo.CreateUser(ctx, req)
}

func (u *useCase) UpdateUser(ctx context.Context, req *auth_v1.UpdateRequest) error {
	return u.repo.UpdateUser(ctx, req)
}

func (u *useCase) DeleteUser(ctx context.Context, id int64) error {
	return u.repo.DeleteUser(ctx, id)
}

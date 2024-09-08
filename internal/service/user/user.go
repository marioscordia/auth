package user

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/marioscordia/auth/internal/client/cache"
	"github.com/marioscordia/auth/internal/model"
	repo "github.com/marioscordia/auth/internal/repository"
	"github.com/marioscordia/auth/internal/service"
)

// New is the function that returns Service object
func New(repo repo.Repository, cache cache.Cache) service.Service {
	return &serve{
		repo:  repo,
		cache: cache,
	}
}

type serve struct {
	repo  repo.Repository
	cache cache.Cache
}

func (u *serve) GetUser(ctx context.Context, id int64) (*model.User, error) {
	user := u.cache.Get(ctx, id)
	if user != nil {
		return user, nil
	}

	user, err := u.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	u.cache.Save(ctx, user)

	return user, nil
}

func (u *serve) CreateUser(ctx context.Context, user *model.UserCreate) (int64, error) {
	user.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))

	id, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	us, err := u.repo.GetUser(ctx, id)
	if err != nil {
		return 0, err
	}

	u.cache.Save(ctx, us)

	return id, nil
}

func (u *serve) UpdateUser(ctx context.Context, user *model.UserUpdate) error {
	if err := u.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	u.cache.Update(ctx, user)

	return nil
}

func (u *serve) DeleteUser(ctx context.Context, id int64) error {
	if err := u.repo.DeleteUser(ctx, id); err != nil {
		return err
	}

	u.cache.Delete(ctx, id)

	return nil
}

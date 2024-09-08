package cache

import (
	"context"

	"github.com/marioscordia/auth/internal/model"
)

type Cache interface {
	Save(ctx context.Context, user *model.User)
	Get(ctx context.Context, userId int64) *model.User
	Delete(ctx context.Context, userId int64)
	Update(ctx context.Context, update *model.UserUpdate)
}

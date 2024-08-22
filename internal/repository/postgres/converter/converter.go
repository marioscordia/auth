package converter

import (
	"github.com/marioscordia/auth/internal/model"
	modelRepo "github.com/marioscordia/auth/internal/repository/postgres/model"
)

// ToUserFromRepo is the method that converts UserDB object to User object, which is used at Handler and Service layers
func ToUserFromRepo(u *modelRepo.UserDB) *model.User {
	return &model.User{
		ID:        u.ID,
		Name:      u.Username,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

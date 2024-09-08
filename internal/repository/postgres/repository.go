package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/marioscordia/auth/internal/model"
	"github.com/marioscordia/auth/internal/repository/postgres/converter"
	modelRepo "github.com/marioscordia/auth/internal/repository/postgres/model"

	repo "github.com/marioscordia/auth/internal/repository"
	"github.com/marioscordia/auth/pkg/auth_v1"
)

// New is a function that returns Repository object
func New(_ context.Context, db *sqlx.DB) (repo.Repository, error) {
	return &repository{
		db: db,
	}, nil
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) CreateUser(ctx context.Context, user *model.UserCreate) (int64, error) {
	var exists bool

	checkName, err := r.db.PreparexContext(ctx, `select exists(select 1 from users where username=$1)`)
	if err != nil {
		return 0, err
	}

	if err = checkName.GetContext(ctx, &exists, user.Name); err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("user with such username already exists")
	}

	checkEmail, err := r.db.PreparexContext(ctx, `select exists(select 1 from users where email=$1)`)
	if err != nil {
		return 0, err
	}

	if err = checkEmail.GetContext(ctx, &exists, user.Email); err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("user with such email already exists")
	}

	role := auth_v1.Role_name[int32(auth_v1.Role_USER)]

	var id int64
	createUser, err := r.db.PreparexContext(ctx, `insert into users (username, email, role, password)
													values ($1, $2, $3, $4) returning id`)

	if err != nil {
		return 0, err
	}

	if err := createUser.GetContext(ctx, &id, user.Name, user.Email, role, user.Password); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	var u modelRepo.UserDB

	getUser, err := r.db.PreparexContext(ctx, `select id, username, email, role, created_at, updated_at from users where id=$1 and deleted_at is null`)
	if err != nil {
		return nil, err
	}

	if err := getUser.GetContext(ctx, &u, id); err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&u), nil
}

func (r *repository) UpdateUser(ctx context.Context, user *model.UserUpdate) error {
	q := "update users set "

	var exists bool

	if user.Name != "" {
		checkName, err := r.db.PreparexContext(ctx, `select exists(select 1 from users where username=$1)`)
		if err != nil {
			return err
		}

		if err := checkName.GetContext(ctx, &exists, user.Name); err != nil {
			return err
		}
		if exists {
			return errors.New("user with such username already exists")
		}

		q += fmt.Sprintf("username='%s'", user.Name)
	}

	if user.Email != "" {
		checkEmail, err := r.db.PreparexContext(ctx, `select exists(select 1 from users where email=$1)`)
		if err != nil {
			return err
		}

		if err := checkEmail.GetContext(ctx, &exists, user.Email); err != nil {
			return err
		}
		if exists {
			return errors.New("user with such email already exists")
		}

		if user.Name != "" {
			q += ", "
		}
		q += fmt.Sprintf("email='%s'", user.Email)
	}

	q += fmt.Sprintf(" where id=%d", user.ID)

	_, err := r.db.ExecContext(ctx, q)

	return err
}

func (r *repository) DeleteUser(ctx context.Context, id int64) error {
	deleteUser, err := r.db.PreparexContext(ctx, `update users set deleted_at=now(), status='inactive' where id=$1`)
	if err != nil {
		return err
	}

	_, err = deleteUser.ExecContext(ctx, id)

	return err
}

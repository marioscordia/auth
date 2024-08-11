package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/marioscordia/auth/internal/model"
	"github.com/marioscordia/auth/internal/repository/postgres/converter"
	modelRepo "github.com/marioscordia/auth/internal/repository/postgres/model"

	repo "github.com/marioscordia/auth/internal/repository"
	"github.com/marioscordia/auth/pkg/auth_v1"
)

// New is a function that returns Repository object
func New(ctx context.Context, db *sqlx.DB) (repo.Repository, error) {
	stmtCreateUser, err := db.PreparexContext(ctx, `insert into users (username, email, role, password, created_at, updated_at)
              									    values ($1, $2, $3, $4, $5, $6) returning id`)

	if err != nil {
		return nil, err
	}

	stmtGetUser, err := db.PreparexContext(ctx, `select id, username, email, role, created_at, updated_at from users where id=$1`)
	if err != nil {
		return nil, err
	}

	stmtDeleteUser, err := db.PreparexContext(ctx, `update users set deleted_at=now(), status='inactive' where id=$1`)
	if err != nil {
		return nil, err
	}

	stmtCheckByName, err := db.PreparexContext(ctx, `select exists(select 1 from users where username=$1)`)
	if err != nil {
		return nil, err
	}

	stmtCheckByEmail, err := db.PreparexContext(ctx, `select exists(select 1 from users where email=$1)`)
	if err != nil {
		return nil, err
	}

	return &repository{
		db:               db,
		stmtCreateUser:   stmtCreateUser,
		stmtGetUser:      stmtGetUser,
		stmtDeleteUser:   stmtDeleteUser,
		stmtCheckByName:  stmtCheckByName,
		stmtCheckByEmail: stmtCheckByEmail,
	}, nil
}

type repository struct {
	db               *sqlx.DB
	stmtCreateUser   *sqlx.Stmt
	stmtGetUser      *sqlx.Stmt
	stmtDeleteUser   *sqlx.Stmt
	stmtCheckByName  *sqlx.Stmt
	stmtCheckByEmail *sqlx.Stmt
}

func (r *repository) CreateUser(ctx context.Context, user *model.User, password string) (int64, error) {
	var exists bool

	if err := r.stmtCheckByName.GetContext(ctx, &exists, user.Name); err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("user with such username already exists")
	}

	if err := r.stmtCheckByEmail.GetContext(ctx, &exists, user.Email); err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("user with such email already exists")
	}

	t := time.Now()
	role := auth_v1.Role_name[int32(auth_v1.Role_USER)]

	var id int64
	if err := r.stmtCreateUser.GetContext(ctx, &id, user.Name, user.Email, role, password, t, t); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	var u modelRepo.UserDB

	if err := r.stmtGetUser.GetContext(ctx, &u, id); err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&u), nil
}

func (r *repository) UpdateUser(ctx context.Context, user *model.User) error {
	q := "update users set "

	var exists bool

	if user.Name != "" {
		if err := r.stmtCheckByName.GetContext(ctx, &exists, user.Name); err != nil {
			return err
		}
		if exists {
			return errors.New("user with such username already exists")
		}

		q += fmt.Sprintf("username='%s'", user.Name)
	}

	if user.Email != "" {
		if err := r.stmtCheckByEmail.GetContext(ctx, &exists, user.Email); err != nil {
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
	_, err := r.stmtDeleteUser.ExecContext(ctx, id)

	return err
}

// type db struct {
// 	ID        int64     `db:"id"`
// 	Username  string    `db:"username"`
// 	Email     string    `db:"email"`
// 	Role      string    `db:"role"`
// 	Status    string    `db:"status"`
// 	CreatedAt time.Time `db:"created_at"`
// 	UpdatedAt time.Time `db:"updated_at"`
// }

// func (u *db) toDto() *model.User {
// 	return &model.User{
// 		ID:        u.ID,
// 		Name:      u.Username,
// 		Email:     u.Email,
// 		Role:      u.Role,
// 		CreatedAt: u.CreatedAt,
// 		UpdatedAt: u.UpdatedAt,
// 	}

// 	// role := auth_v1.Role(auth_v1.Role_value[u.Role])

// 	// return &auth_v1.GetResponse{
// 	// 	Id:        u.ID,
// 	// 	Name:      u.Username,
// 	// 	Email:     u.Email,
// 	// 	Role:      role,
// 	// 	CreatedAt: timestamppb.New(u.CreatedAt),
// 	// 	UpdatedAt: timestamppb.New(u.UpdatedAt),
// 	// }
// }

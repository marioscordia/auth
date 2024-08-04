package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/marioscordia/auth/pkg/auth_v1"
	repo "github.com/marioscordia/auth/repository"
)

// New is ...
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

func (r *repository) CreateUser(ctx context.Context, u *auth_v1.CreateRequest) (int64, error) {
	var exists bool

	if err := r.stmtCheckByName.GetContext(ctx, &exists, u.Name); err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("user with such username already exists")
	}

	if err := r.stmtCheckByEmail.GetContext(ctx, &exists, u.Email); err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("user with such email already exists")
	}

	t := time.Now()
	role := auth_v1.Role_name[int32(auth_v1.Role_USER)]

	var id int64
	if err := r.stmtCreateUser.GetContext(ctx, &id, u.Name, u.Email, role, u.Password, t, t); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) GetUser(ctx context.Context, id int64) (*auth_v1.GetResponse, error) {
	var u db

	if err := r.stmtGetUser.GetContext(ctx, &u, id); err != nil {
		return nil, err
	}

	return u.toDto(), nil
}

func (r *repository) UpdateUser(ctx context.Context, u *auth_v1.UpdateRequest) error {
	q := "update users set "

	var exists bool

	if u.Name.Value != "" {
		if err := r.stmtCheckByName.GetContext(ctx, &exists, u.Name.Value); err != nil {
			return err
		}
		if exists {
			return errors.New("user with such username already exists")
		}

		q += fmt.Sprintf("username='%s'", u.Name.Value)
	}

	if u.Email.Value != "" {
		if err := r.stmtCheckByEmail.GetContext(ctx, &exists, u.Email.Value); err != nil {
			return err
		}
		if exists {
			return errors.New("user with such email already exists")
		}

		if u.Name.Value != "" {
			q += ", "
		}
		q += fmt.Sprintf("email='%s'", u.Email.Value)
	}

	q += fmt.Sprintf(" where id=%d", u.Id)

	_, err := r.db.ExecContext(ctx, q)

	return err
}

func (r *repository) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.stmtDeleteUser.ExecContext(ctx, id)

	return err
}

type db struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *db) toDto() *auth_v1.GetResponse {
	role := auth_v1.Role(auth_v1.Role_value[u.Role])

	return &auth_v1.GetResponse{
		Id:        u.ID,
		Name:      u.Username,
		Email:     u.Email,
		Role:      role,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

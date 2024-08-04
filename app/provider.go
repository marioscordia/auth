package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // to initialize connection
	"github.com/pressly/goose/v3"

	"github.com/marioscordia/auth/closer"
	"github.com/marioscordia/auth/delivery/grpc"
	"github.com/marioscordia/auth/facility"
	repo "github.com/marioscordia/auth/repository"
	"github.com/marioscordia/auth/repository/postgres"
	"github.com/marioscordia/auth/service"
	"github.com/marioscordia/auth/service/user"
)

const (
	dbPostgresDriverName   = "postgres"
	migrationsPostgresPath = "db/migrations"
)

type provider struct {
	config *facility.Config

	db *sqlx.DB

	repo repo.Repository

	service service.Service

	handler *grpc.Handler
}

func newProvider() *provider {
	return &provider{}
}

func (p *provider) NewDB() *sqlx.DB {
	if p.db == nil {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			p.config.PostgresHost, p.config.PostgresPort, p.config.PostgresUser, p.config.PostgresPassword, p.config.PostgresDb, p.config.PostgresSslMode)

		db, err := sql.Open(dbPostgresDriverName, psqlInfo)
		if err != nil {
			log.Fatalf("failed to connect to postgres db: %v", err)
		}

		dbx := sqlx.NewDb(db, dbPostgresDriverName)

		if err = dbx.Ping(); err != nil {
			log.Fatalf("failed to verify connection: %v", err)

		}

		if p.config.PostgresMigrate {
			if err := goose.SetDialect("postgres"); err != nil {
				log.Fatalf("failed to set postgres dialect for goose: %v", err)
			}
			if err := goose.Up(db, migrationsPostgresPath); err != nil && !errors.Is(err, goose.ErrAlreadyApplied) {
				log.Fatalf("failed to apply migrations: %v", err)
			}
		}

		p.db = dbx

		closer.Add(dbx.Close)
	}

	return p.db
}

func (p *provider) UserRepository(ctx context.Context) repo.Repository {
	if p.repo == nil {
		repo, err := postgres.New(ctx, p.NewDB())
		if err != nil {
			log.Fatalf("failed to initialize chat repository: %v", err)
		}

		p.repo = repo
	}

	return p.repo
}

func (p *provider) UserService(ctx context.Context) service.Service {
	if p.service == nil {
		p.service = user.New(p.UserRepository(ctx))
	}

	return p.service
}

func (p *provider) UserHandler(ctx context.Context) *grpc.Handler {
	if p.handler == nil {
		p.handler = grpc.New(p.UserService(ctx))
	}

	return p.handler
}

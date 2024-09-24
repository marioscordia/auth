package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // to initialize connection
	"github.com/pressly/goose/v3"

	"github.com/marioscordia/auth/internal/api"
	"github.com/marioscordia/auth/internal/client/cache"
	redisgo "github.com/marioscordia/auth/internal/client/cache/redis"
	"github.com/marioscordia/auth/internal/closer"
	"github.com/marioscordia/auth/internal/config"
	repo "github.com/marioscordia/auth/internal/repository"
	"github.com/marioscordia/auth/internal/repository/postgres"
	"github.com/marioscordia/auth/internal/service"
	"github.com/marioscordia/auth/internal/service/user"
)

const (
	dbPostgresDriverName   = "postgres"
	migrationsPostgresPath = "db/migrations"
)

type provider struct {
	config *config.Config

	db *sqlx.DB

	redisConn redis.Conn

	repository repo.Repository

	userCache cache.Cache

	service service.Service

	handler *api.Handler
}

func newProvider() *provider {
	return &provider{}
}

// NewDB initializes the connection to the PostgreSQL
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
			if err = goose.SetDialect("postgres"); err != nil {
				log.Fatalf("failed to set postgres dialect for goose: %v", err)
			}

			if err = goose.Up(db, migrationsPostgresPath); err != nil {
				log.Fatalf("failed to apply migrations: %v", err)
			}
		}

		p.db = dbx

		closer.Add(dbx.Close)
	}

	return p.db
}

// UserRepository initializes the repository interface using database connection
func (p *provider) UserRepository(ctx context.Context) repo.Repository {
	if p.repository == nil {
		repo, err := postgres.New(ctx, p.NewDB())
		if err != nil {
			log.Fatalf("failed to initialize chat repository: %v", err)
		}

		p.repository = repo
	}

	return p.repository
}

// NewRedisConn initializes Redis connection
func (p *provider) NewRedisConn(ctx context.Context) redis.Conn {
	if p.redisConn == nil {
		connInfo := fmt.Sprintf("%s:%d", p.config.RedisHost, p.config.RedisPort)

		conn, err := redis.DialContext(
			ctx,
			"tcp",
			connInfo,
		)

		if err != nil {
			log.Fatalf("failed to connect to redis: %v", err)
		}

		p.redisConn = conn

		closer.Add(conn.Close)
	}

	return p.redisConn
}

// NewUserCache initializes the cache interface using Redis connection
func (p *provider) UserCache(ctx context.Context) cache.Cache {
	if p.userCache == nil {
		p.userCache = redisgo.NewUserCache(p.NewRedisConn(ctx))
	}

	return p.userCache
}

// UserService initializes the service interface using UserRepository and UserCache interfaces
func (p *provider) UserService(ctx context.Context) service.Service {
	if p.service == nil {
		p.service = user.New(p.UserRepository(ctx), p.UserCache(ctx))
	}

	return p.service
}

// UserHandler initializes the API handler using UserService interface
func (p *provider) UserHandler(ctx context.Context) *api.Handler {
	if p.handler == nil {
		p.handler = api.New(p.UserService(ctx))
	}

	return p.handler
}

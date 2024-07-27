package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"

	"github.com/marioscordia/auth/app"
	"github.com/marioscordia/auth/facility"
)

const (
	dbPostgresDriverName   = "postgres"
	migrationsPostgresPath = "db/migrations"
)

func main() {
	cfg, err := facility.NewConfig()
	if err != nil {
		log.Panicf("failed to create config: %v", err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDb, cfg.PostgresSslMode)

	db, err := sql.Open(dbPostgresDriverName, psqlInfo)
	if err != nil {
		log.Panicf("failed to connect to postgres db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Panicf("error closing the db: %v", err)
		}
	}()

	dbx := sqlx.NewDb(db, dbPostgresDriverName)

	if cfg.PostgresMigrate {
		if err := goose.SetDialect("postgres"); err != nil {
			log.Panicf("failed to set postgres dialect for goose: %v", err)
		}
		if err := goose.Up(db, migrationsPostgresPath); err != nil && !errors.Is(err, goose.ErrAlreadyApplied) {
			log.Panicf("failed to apply migrations: %v", err)
		}
	}

	server := grpc.NewServer()

	startingCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(startingCtx, dbx, server, cfg); err != nil {
		log.Panicf("failed to run app: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	server.GracefulStop()
}

package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	user "github.com/marioscordia/auth"
	authGrpc "github.com/marioscordia/auth/delivery/grpc"
	"github.com/marioscordia/auth/facility"
	"github.com/marioscordia/auth/pkg/auth_v1"
	"github.com/marioscordia/auth/repository/postgres"
)

// Run is ...
func Run(ctx context.Context, postgresDb *sqlx.DB, server *grpc.Server, config *facility.Config) error {
	repo, err := postgres.New(ctx, postgresDb)
	if err != nil {
		return err
	}

	useCase := user.New(repo)

	handler := authGrpc.New(useCase)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcPort))
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	defer func() {
		if err = lis.Close(); err != nil {
			log.Panicf("error closing the listener: %v", err)
		}
	}()

	reflection.Register(server)
	auth_v1.RegisterAuthV1Server(server, handler)

	log.Printf("server listening at %v", lis.Addr())

	if err = server.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}

	return nil
}

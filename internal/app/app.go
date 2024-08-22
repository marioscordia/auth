package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/marioscordia/auth/internal/closer"
	"github.com/marioscordia/auth/internal/config"

	"github.com/marioscordia/auth/pkg/auth_v1"
)

// App is an object with initialzing and starting methods
type App struct {
	provider   *provider
	grpcServer *grpc.Server
}

// NewApp is function that returns App object
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initProvider,
		a.initConfig,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// Run is a method that starts the application
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) initProvider(_ context.Context) error {
	a.provider = newProvider()
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	a.provider.config = cfg

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	auth_v1.RegisterAuthV1Server(a.grpcServer, a.provider.UserHandler(ctx))

	closer.Add(a.gracefulStop)

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %d", a.provider.config.GrpcPort)

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", a.provider.config.GrpcPort))
	if err != nil {
		return err
	}

	closer.Add(list.Close)

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) gracefulStop() error {
	a.grpcServer.GracefulStop()
	return nil
}

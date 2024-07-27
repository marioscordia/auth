package main

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/marioscordia/auth/pkg/auth_v1"
)

type server struct {
	auth_v1.UnimplementedAuthV1Server
}

func (s *server) Get(_ context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	log.Printf("Getting User with id: %d\n", req.GetId())

	return &auth_v1.GetResponse{
		Id:        req.Id,
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      0,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *server) Create(_ context.Context, _ *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	log.Printf("Creating User...")

	userID := gofakeit.Uint64()

	log.Printf("User created successfully with id: %d", userID)
	return &auth_v1.CreateResponse{
		Id: int64(userID),
	}, nil
}

func (s *server) Update(_ context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Updating User with id %d", req.Id)
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Deleting User with id %d", req.Id)
	return &emptypb.Empty{}, nil
}

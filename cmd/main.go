package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"github.com/marioscordia/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

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

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	auth_v1.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

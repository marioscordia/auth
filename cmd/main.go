package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/marioscordia/auth/pkg/auth_v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	auth_v1.RegisterAuthV1Server(s, &server{})

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("server listening at %v", lis.Addr())

		if err = s.Serve(lis); err != nil {
			log.Panicf("failed to serve: %v", err)
		}
	}()

	<-signalChan
	log.Println("received shutdown signal")

	s.GracefulStop()

	log.Println("server shutdown complete")
}

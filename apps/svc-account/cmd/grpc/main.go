package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/buemura/minibank/packages/pb"
	"github.com/buemura/minibank/svc-account/config"
	"github.com/buemura/minibank/svc-account/internal/infra/database"
	"github.com/buemura/minibank/svc-account/internal/infra/handler"
	"google.golang.org/grpc"
)

func init() {
	config.LoadEnv()
	database.Connect()
}

func main() {
	port := ":" + config.GRPC_PORT
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot create grpc listener: %s", err)
	}

	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, &handler.AccountHandler{})

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to server grpc: %s", err)
		}
	}()

	log.Println("gRPC Server running at", port, "...")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)
	<-stop

	log.Println("Stopping gRPC Server...")
	s.GracefulStop()
	log.Println("gRPC Server stopped")
}

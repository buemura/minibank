package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/buemura/minibank/packages/gen/protos"
	"github.com/buemura/minibank/svc-transaction/config"
	"github.com/buemura/minibank/svc-transaction/internal/infra/database"
	"github.com/buemura/minibank/svc-transaction/internal/infra/event"
	"github.com/buemura/minibank/svc-transaction/internal/infra/handler"
	"google.golang.org/grpc"
)

func init() {
	config.LoadEnv()
	database.Connect()
}

func main() {
	// Start queue consumer in a separate goroutine to handle incoming events.
	go event.QueueConsumer()

	port := ":" + config.GRPC_PORT
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot create grpc listener: %s", err)
	}

	s := grpc.NewServer()
	protos.RegisterTransactionServiceServer(s, &handler.TransactionHandler{})

	// Run gRPC server in a separate goroutine.
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

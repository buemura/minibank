package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buemura/minibank/api-gateway/config"
	"github.com/buemura/minibank/api-gateway/internal/infra/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	config.LoadEnv()
}

func main() {
	server := echo.New()
	server.Use(middleware.CORS())

	handler.SetupRoutes(server)

	port := ":" + config.PORT

	go func() {
		if err := server.Start(port); err != nil && http.ErrServerClosed != err {
			panic(err)
		}
	}()

	log.Println("HTTP Server running at", port, "...")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Stopping HTTP Server...")

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Println("HTTP Server stopped")
}

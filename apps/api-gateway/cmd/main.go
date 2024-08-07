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
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	})
)

func init() {
	config.LoadEnv()
	prometheus.MustRegister(requestsTotal)
}

func main() {
	server := echo.New()
	server.Use(middleware.CORS())
	server.Use(echoprometheus.NewMiddleware("myapp"))
	server.GET("/metrics", echoprometheus.NewHandler())

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

package handlers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const healthCheckTimeout = 3 * time.Second

type ServiceHealth struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Latency string `json:"latency"`
}

type HealthHandler struct {
	authConn        *grpc.ClientConn
	accountConn     *grpc.ClientConn
	transactionConn *grpc.ClientConn
}

func NewHealthHandler(authConn, accountConn, transactionConn *grpc.ClientConn) *HealthHandler {
	return &HealthHandler{
		authConn:        authConn,
		accountConn:     accountConn,
		transactionConn: transactionConn,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	services := h.checkAllServices(c.Request.Context())

	overallStatus := "healthy"
	for _, svc := range services {
		if svc.Status != "SERVING" {
			overallStatus = "degraded"
			break
		}
	}

	statusCode := http.StatusOK
	if overallStatus == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":   overallStatus,
		"services": services,
	})
}

func (h *HealthHandler) checkAllServices(ctx context.Context) []ServiceHealth {
	type result struct {
		index  int
		health ServiceHealth
	}

	conns := []struct {
		name string
		conn *grpc.ClientConn
	}{
		{"svc-auth", h.authConn},
		{"svc-account", h.accountConn},
		{"svc-transaction", h.transactionConn},
	}

	results := make([]ServiceHealth, len(conns))
	var wg sync.WaitGroup
	ch := make(chan result, len(conns))

	for i, svc := range conns {
		wg.Add(1)
		go func(idx int, name string, conn *grpc.ClientConn) {
			defer wg.Done()
			ch <- result{index: idx, health: checkService(ctx, name, conn)}
		}(i, svc.name, svc.conn)
	}

	wg.Wait()
	close(ch)

	for r := range ch {
		results[r.index] = r.health
	}
	return results
}

func checkService(ctx context.Context, name string, conn *grpc.ClientConn) ServiceHealth {
	client := healthpb.NewHealthClient(conn)
	checkCtx, cancel := context.WithTimeout(ctx, healthCheckTimeout)
	defer cancel()

	start := time.Now()
	resp, err := client.Check(checkCtx, &healthpb.HealthCheckRequest{
		Service: "",
	})
	latency := time.Since(start)

	if err != nil {
		return ServiceHealth{
			Name:    name,
			Status:  "UNAVAILABLE",
			Latency: latency.String(),
		}
	}

	return ServiceHealth{
		Name:    name,
		Status:  resp.Status.String(),
		Latency: latency.String(),
	}
}

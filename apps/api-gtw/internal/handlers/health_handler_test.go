package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func startTestHealthServer(t *testing.T, serving bool) *grpc.ClientConn {
	t.Helper()
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	s := grpc.NewServer()
	hs := health.NewServer()
	if serving {
		hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	} else {
		hs.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	}
	healthpb.RegisterHealthServer(s, hs)

	go s.Serve(lis)
	t.Cleanup(func() { s.Stop() })

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return conn
}

func TestHealthCheck_AllHealthy(t *testing.T) {
	authConn := startTestHealthServer(t, true)
	accountConn := startTestHealthServer(t, true)
	transactionConn := startTestHealthServer(t, true)

	handler := NewHealthHandler(authConn, accountConn, transactionConn)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/health", nil)

	handler.HealthCheck(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "healthy", resp["status"])

	services := resp["services"].([]interface{})
	assert.Len(t, services, 3)
	for _, s := range services {
		svc := s.(map[string]interface{})
		assert.Equal(t, "SERVING", svc["status"])
	}
}

func TestHealthCheck_Degraded(t *testing.T) {
	authConn := startTestHealthServer(t, true)
	accountConn := startTestHealthServer(t, false)
	transactionConn := startTestHealthServer(t, true)

	handler := NewHealthHandler(authConn, accountConn, transactionConn)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/health", nil)

	handler.HealthCheck(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "degraded", resp["status"])

	services := resp["services"].([]interface{})
	for _, s := range services {
		svc := s.(map[string]interface{})
		if svc["name"] == "svc-account" {
			assert.Equal(t, "NOT_SERVING", svc["status"])
		}
	}
}

func TestHealthCheck_ServiceUnavailable(t *testing.T) {
	authConn := startTestHealthServer(t, true)
	// Connect to a port that's not listening
	badConn, err := grpc.NewClient("localhost:1",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { badConn.Close() })
	transactionConn := startTestHealthServer(t, true)

	handler := NewHealthHandler(authConn, badConn, transactionConn)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/health", nil)

	handler.HealthCheck(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "degraded", resp["status"])

	services := resp["services"].([]interface{})
	for _, s := range services {
		svc := s.(map[string]interface{})
		if svc["name"] == "svc-account" {
			assert.Equal(t, "UNAVAILABLE", svc["status"])
		}
	}
}

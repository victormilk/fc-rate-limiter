package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/victormilk/fc-rate-limiter/limiter"
)

func setupRedis(t *testing.T) (*redis.Client, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, err := redisC.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: host + ":" + port.Port(),
	})

	return client, func() {
		client.Close()
		redisC.Terminate(ctx)
	}
}

func TestRateLimiter(t *testing.T) {
	client, teardown := setupRedis(t)
	defer teardown()

	rateLimiter := limiter.NewRedisLimiter(client.Options().Addr, 5, 10, 10*time.Second)

	handler := RateLimiter(rateLimiter, 5, 10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("allow requests under limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()

		for i := 0; i < 5; i++ {
			handler.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("block requests over limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()

		for i := 0; i < 6; i++ {
			handler.ServeHTTP(w, req)
		}
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("allow requests with token under limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("API_KEY", "test-token")
		w := httptest.NewRecorder()

		for i := 0; i < 10; i++ {
			handler.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("block requests with token over limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("API_KEY", "test-token")
		w := httptest.NewRecorder()

		for i := 0; i < 11; i++ {
			handler.ServeHTTP(w, req)
		}
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})
}

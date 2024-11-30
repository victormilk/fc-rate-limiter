package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/victormilk/fc-rate-limiter/limiter"
	"github.com/victormilk/fc-rate-limiter/middleware"
)

func main() {
	r := chi.NewRouter()

	redisHost := os.Getenv("REDIS_HOST")
	ipRate, err := strconv.Atoi(os.Getenv("IP_RATE"))
	if err != nil {
		log.Fatal(err)
	}
	tokenRate, err := strconv.Atoi(os.Getenv("TOKEN_RATE"))
	if err != nil {
		log.Fatal(err)
	}
	blockDuration, err := strconv.Atoi(os.Getenv("BLOCK_DURATION"))
	if err != nil {
		log.Fatal(err)
	}
	limiter := limiter.NewRedisLimiter(redisHost, ipRate, tokenRate, time.Duration(blockDuration))
	r.Use(middleware.RateLimiter(limiter, 2, 2))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	log.Printf("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

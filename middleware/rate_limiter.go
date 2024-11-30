package middleware

import (
	"net/http"
	"strings"

	"github.com/victormilk/fc-rate-limiter/limiter"
)

const (
	blockMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func RateLimiter(l limiter.Limiter, ipLimit, tokenLimit int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ip := strings.Split(r.RemoteAddr, ":")[0]
			token := r.Header.Get("API_KEY")

			var key string
			if token != "" {
				key = token
			} else {
				key = ip
			}

			blocked, err := l.IsBlocked(ctx, key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if blocked {
				http.Error(w, blockMessage, http.StatusTooManyRequests)
				return
			}

			allowed, err := l.Allow(ctx, key, ipLimit)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !allowed {
				l.Block(ctx, key)
				http.Error(w, blockMessage, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

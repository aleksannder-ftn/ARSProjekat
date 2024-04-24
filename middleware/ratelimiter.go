package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mux      sync.Mutex
	interval time.Duration
	limit    int
	counters map[string]int
}

func NewRateLimiter(interval time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		interval: interval,
		limit:    limit,
		counters: make(map[string]int),
	}
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		rl.mux.Lock()
		defer rl.mux.Unlock()

		count, ok := rl.counters[ip]
		if !ok {
			count = 1
		} else {
			count++
		}

		if count > rl.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		rl.counters[ip] = count
		time.AfterFunc(rl.interval, func() {
			rl.mux.Lock()
			defer rl.mux.Unlock()
			delete(rl.counters, ip)
		})

		next.ServeHTTP(w, r)
	})
}

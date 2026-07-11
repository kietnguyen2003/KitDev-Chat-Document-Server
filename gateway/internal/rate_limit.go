package gateway

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	Count   int
	ResetAt time.Time
}

type RateLimiter struct {
	limit         int
	window        time.Duration
	cleanupWindow time.Duration
	mu            sync.Mutex
	clients       map[string]*rateLimitEntry
}

func NewRateLimiter(limit int, window, cleanupWindow time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		limit:         limit,
		window:        window,
		cleanupWindow: cleanupWindow,
		clients:       make(map[string]*rateLimitEntry),
	}

	go limiter.cleanupLoop()

	return limiter
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		now := time.Now()

		allowed, remaining, resetAt := rl.allow(key, now)

		c.Header("X-RateLimit-Limit", itoa(rl.limit))
		c.Header("X-RateLimit-Remaining", itoa(remaining))
		c.Header("X-RateLimit-Reset", itoa(int(resetAt.Unix())))

		if allowed {
			c.Next()
			return
		}

		retryAfter := int(time.Until(resetAt).Seconds())
		if retryAfter < 1 {
			retryAfter = 1
		}

		c.Header("Retry-After", itoa(retryAfter))
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"code": http.StatusTooManyRequests,
			"msg":  "Rate limit exceeded",
			"data": nil,
		})
	}
}

func (rl *RateLimiter) allow(key string, now time.Time) (bool, int, time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.clients[key]
	if !exists || now.After(entry.ResetAt) {
		entry = &rateLimitEntry{
			Count:   0,
			ResetAt: now.Add(rl.window),
		}
		rl.clients[key] = entry
	}

	if entry.Count >= rl.limit {
		return false, 0, entry.ResetAt
	}

	entry.Count++
	remaining := rl.limit - entry.Count

	return true, remaining, entry.ResetAt
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanupWindow)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanupExpired(time.Now())
	}
}

func (rl *RateLimiter) cleanupExpired(now time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for key, entry := range rl.clients {
		if now.After(entry.ResetAt) {
			delete(rl.clients, key)
		}
	}
}

func itoa(value int) string {
	return fmtInt(int64(value))
}

func fmtInt(value int64) string {
	if value == 0 {
		return "0"
	}

	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}

	buf := [20]byte{}
	i := len(buf)

	for value > 0 {
		i--
		buf[i] = byte('0' + value%10)
		value /= 10
	}

	return sign + string(buf[i:])
}

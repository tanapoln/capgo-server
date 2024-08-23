package ratelimit

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

var limiterSet = cache.New(5*time.Minute, 10*time.Minute)

type LimiterKeyFunc func(*gin.Context) string
type LimiterFactory func(*gin.Context) (*rate.Limiter, time.Duration)
type AbortFunc func(*gin.Context)

func NewRateLimiter(key LimiterKeyFunc, factory LimiterFactory, abort AbortFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		k := key(c)
		limiter, ok := limiterSet.Get(k)
		if !ok {
			limiter, expire := factory(c)
			limiterSet.Set(k, limiter, expire)
		}

		ok = limiter.(*rate.Limiter).Allow()
		if !ok {
			abort(c)
			return
		}
		c.Next()
	}
}

func DefaultAbort(c *gin.Context) {
	c.AbortWithStatus(http.StatusTooManyRequests)
}

func CreateLimiterFactory(lim rate.Limit, count int) LimiterFactory {
	return func(ctx *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(lim, count), 1 * time.Hour
	}
}

func KeyByIPAddress(c *gin.Context) string {
	return c.ClientIP()
}

package apikey

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewApiKeyMiddleware(headerKey string, keys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := strings.TrimSpace(c.GetHeader(headerKey))
		if apiKey == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !slices.Contains(keys, apiKey) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

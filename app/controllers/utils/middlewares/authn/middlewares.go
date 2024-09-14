package authn

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/tanapoln/capgo-server/config"
	"golang.org/x/oauth2"
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

func NewOAuthMiddleware(headerKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		oauthToken := strings.TrimSpace(c.GetHeader(headerKey))
		if oauthToken == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(oauthToken, " ")
		if len(parts) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authType, authToken := parts[0], parts[1]
		if authType != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		provider, err := oidc.NewProvider(c.Request.Context(), config.Get().OAuthIssuer)
		if err != nil {
			slog.Error("Error creating OAuth provider", "error", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		_, err = provider.UserInfo(c.Request.Context(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: authToken}))
		if err != nil {
			slog.Info("Error getting user info", "error", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func MultiAuthMiddleware(middlewares map[string]gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, middleware := range middlewares {
			if c.GetHeader(key) != "" {
				middleware(c)
				return
			}
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	capgoCtrl "github.com/tanapoln/capgo-server/app/controllers/capgo"
	mgmtCtrl "github.com/tanapoln/capgo-server/app/controllers/mgmt"
	"github.com/tanapoln/capgo-server/app/controllers/utils/middlewares/authn"
	"github.com/tanapoln/capgo-server/app/controllers/utils/middlewares/httpstats"
	"github.com/tanapoln/capgo-server/app/controllers/utils/middlewares/ratelimit"
	"github.com/tanapoln/capgo-server/app/controllers/utils/middlewares/spa"
	"github.com/tanapoln/capgo-server/config"
	"golang.org/x/time/rate"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(config.Get().TrustedProxies)

	router.Use(httpstats.NewMiddleware())

	capgo := router.Group("")
	{
		updateLimit := ratelimit.NewRateLimiter(
			ratelimit.KeyByIPAddress,
			ratelimit.CreateLimiterFactory(rate.Every(1*time.Minute), config.Get().LimitRequestPerMinute),
			ratelimit.DefaultAbort,
		)

		ctrl := capgoCtrl.NewCapgoController()
		capgo.POST("/updates", updateLimit, ctrl.Updates)
		capgo.POST("/stats", ctrl.Stats)
		capgo.POST("/channel_self", ctrl.RegisterChannel)
		capgo.DELETE("/channel_self", ctrl.UnregisterChannel)
	}

	router.GET("/_healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return router
}

func InitMgmtRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(config.Get().TrustedProxies)

	router.Use(httpstats.NewMiddleware())

	mgmt := router.Group("/api/v1/")
	{
		mgmt.Use(authn.MultiAuthMiddleware(map[string]gin.HandlerFunc{
			"x-api-key":     authn.NewApiKeyMiddleware("x-api-key", config.Get().ManagementAPITokens),
			"Authorization": authn.NewOAuthMiddleware("Authorization"),
		}))

		ctrl := mgmtCtrl.NewCapgoManagementController()
		mgmt.GET("/bundles.list", ctrl.ListAllBundles)
		mgmt.POST("/bundles.upload", ctrl.UploadBundle)

		mgmt.GET("/releases.list", ctrl.ListAllReleases)
		mgmt.POST("/releases.create", ctrl.CreateRelease)
		mgmt.POST("/releases.update", ctrl.UpdateRelease)
		mgmt.POST("/releases.set-active", ctrl.SetReleaseActiveBundle)
		mgmt.POST("/releases.delete", ctrl.DeleteRelease)
	}

	router.GET("/_healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	public := router.Group("/apipublic/v1")
	{
		public.GET("/oauth2.config", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"issuer":    config.Get().OAuthIssuer,
				"client_id": config.Get().OAuthClientID,
			})
		})
	}

	spa.Middleware(router, "/ui", "./client/dist", "/index.html")

	return router
}

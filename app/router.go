package app

import (
	"time"

	"github.com/gin-gonic/gin"
	capgoCtrl "github.com/tanapoln/capgo-server/app/controllers/capgo"
	mgmtCtrl "github.com/tanapoln/capgo-server/app/controllers/mgmt"
	"github.com/tanapoln/capgo-server/app/controllers/utils/middlewares/apikey"
	"github.com/tanapoln/capgo-server/app/controllers/utils/middlewares/ratelimit"
	"github.com/tanapoln/capgo-server/config"
	"golang.org/x/time/rate"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(config.Get().TrustedProxies)

	capgo := router.Group("")
	{
		capgo.Use(
			ratelimit.NewRateLimiter(
				ratelimit.KeyByIPAddress,
				ratelimit.CreateLimiterFactory(rate.Every(1*time.Minute), config.Get().LimitRequestPerMinute),
				ratelimit.DefaultAbort,
			),
		)

		ctrl := capgoCtrl.NewCapgoController()
		capgo.POST("/updates", ctrl.Updates)
		capgo.POST("/stats", ctrl.Stats)
		capgo.POST("/channel_self", ctrl.RegisterChannel)
		capgo.DELETE("/channel_self", ctrl.UnregisterChannel)
	}

	mgmt := router.Group("/api/v1/")
	{
		mgmt.Use(apikey.NewApiKeyMiddleware("x-api-key", config.Get().ManagementAPIToken))

		ctrl := mgmtCtrl.NewCapgoManagementController()
		mgmt.GET("/bundles.list", ctrl.ListAllBundles)
		mgmt.POST("/bundles.upload", ctrl.UploadBundle)

		mgmt.GET("/releases.list", ctrl.ListAllReleases)
		mgmt.POST("/releases.create", ctrl.CreateRelease)
		mgmt.POST("/releases.update", ctrl.UpdateRelease)
		mgmt.POST("/releases.set-active", ctrl.SetReleaseActiveBundle)
	}

	return router
}

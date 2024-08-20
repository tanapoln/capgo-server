package app

import (
	"github.com/gin-gonic/gin"
	capgoCtrl "github.com/tanapoln/capgo-server/app/controllers/capgo"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	capgo := router.Group("")
	{
		ctrl := capgoCtrl.NewCapgoController()

		capgo.POST("/updates", ctrl.Updates)
		capgo.POST("/stats", ctrl.Stats)
		capgo.POST("/channel_self", ctrl.RegisterChannel)
		capgo.DELETE("/channel_self", ctrl.UnregisterChannel)
	}

	return router
}

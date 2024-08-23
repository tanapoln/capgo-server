package capgo

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/tanapoln/capgo-server/app/controllers/utils"
	"github.com/tanapoln/capgo-server/app/services"
)

func NewCapgoController() *CapgoController {
	return &CapgoController{
		updateService: &services.UpdateService{},
	}
}

type CapgoController struct {
	updateService *services.UpdateService
}

func (ctrl *CapgoController) Updates(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		var reqBody UpdateRequest
		err := ctx.BindJSON(&reqBody)
		if err != nil {
			return CapgoErrorResponse{
				Error: "invalid request body json",
			}, nil
		}

		if !reqBody.IsValid() {
			return CapgoErrorResponse{
				Error: "invalid request data",
			}, nil
		}

		bundle, err := ctrl.updateService.GetLatest(ctx.Request.Context(), services.GetLatestQuery{
			AppID:       reqBody.AppID,
			Platform:    reqBody.GetPlatform(),
			VersionName: reqBody.VersionName,
			VersionCode: reqBody.VersionCode,
		})
		if err != nil {
			return CapgoErrorResponse{
				Error: err.Error(),
			}, nil
		}

		return UpdateWithNewMinorVersionResponse{
			Version:   bundle.VersionName,
			Checksum:  bundle.CRC,
			URL:       bundle.PublicDownloadURL,
			Signature: bundle.Signature,
		}, nil
	})
}

func (ctrl *CapgoController) Stats(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		slog.Info("Capgo - stats", "body", ctx.Request.Body)
		return gin.H{}, nil
	})
}

func (ctrl *CapgoController) RegisterChannel(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		slog.Info("Capgo - register channel", "body", ctx.Request.Body)
		return gin.H{}, nil
	})
}

func (ctrl *CapgoController) UnregisterChannel(ctx *gin.Context) {
	utils.Handle(ctx, func() (interface{}, error) {
		slog.Info("Capgo - unregister channel", "body", ctx.Request.Body)
		return gin.H{}, nil
	})
}

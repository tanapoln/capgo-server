package capgo

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func NewCapgoController() *CapgoController {
	return &CapgoController{}
}

type CapgoController struct {
}

func (ctrl *CapgoController) Updates(ctx *gin.Context) {
	handle(ctx, func() (interface{}, error) {
		return nil, nil
	})
}

func (ctrl *CapgoController) Stats(ctx *gin.Context) {
	handle(ctx, func() (interface{}, error) {
		return nil, nil
	})
}

func (ctrl *CapgoController) RegisterChannel(ctx *gin.Context) {
	handle(ctx, func() (interface{}, error) {
		return nil, nil
	})
}

func (ctrl *CapgoController) UnregisterChannel(ctx *gin.Context) {
	handle(ctx, func() (interface{}, error) {
		return nil, nil
	})
}

func handle(ctx *gin.Context, fn func() (interface{}, error)) {
	resp, err := fn()
	if err != nil {
		traceId := xid.New().String()
		slog.Error("handler return error. response with HTTP 500", "trace", traceId, "error", err, "path", ctx.Request.RequestURI)
		ctx.JSON(http.StatusInternalServerError, CapgoErrorResponse{
			Error: fmt.Sprintf("Internal Server Error. trace=%s", traceId),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

package utils

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func Handle(ctx *gin.Context, fn func() (interface{}, error)) {
	resp, err := fn()

	if err != nil {
		traceId := xid.New().String()
		slog.Error("handler return error. response with HTTP 500", "trace", traceId, "error", err, "path", ctx.Request.RequestURI)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Internal Server Error. trace=%s", traceId),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

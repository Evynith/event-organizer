package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func inspectError(ctx *gin.Context, err error) {
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func bindJSON(ctx *gin.Context, obj any) {
	err := ctx.BindJSON(&obj)
	inspectError(ctx, err)
}

func bindQuery(ctx *gin.Context, obj any) {
	err := ctx.BindQuery(&obj)
	inspectError(ctx, err)
}

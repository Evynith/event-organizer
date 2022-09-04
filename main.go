package main

import (
	r "main/internal/route"

	"github.com/gin-gonic/gin"
)

func main() {
	route := r.SetupRouter()
	route.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "funciona"})
	})
	route.Run(":8000")
}

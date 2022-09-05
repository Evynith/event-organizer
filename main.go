package main

import (
	"github.com/gin-gonic/gin"

	router "main/internal/route"
)

func main() {
	route := router.SetupRouter()
	route.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Welcome to API event-organizer"})
	})
	route.Run(":8000")
}

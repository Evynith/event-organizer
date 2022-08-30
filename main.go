package main

import (
	"main/internal/controller"
	"net/http"

	authHandler "main/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/login", func(ctx *gin.Context) {
		token := controller.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})
	var typeUser string = "admin"
	public := router.Group("/events")
	public.Use(authHandler.AuthorizeJWT())
	{
		public.GET("", controller.Events)
		public.GET(":id", controller.Event)
		private := public.Group("")
		private.Use(authHandler.OnlyAdmin(typeUser))
		{
			private.POST("", controller.PostEvent)
			private.DELETE(":id", controller.DeleteEvent)
			private.PUT(":id", controller.PutEvent)
		}

		inscription := public.Group("inscription")
		inscription.Use()
		{
			inscription.POST("", controller.PostInscription)
		}
	}

	router.Run("localhost:8080")
}

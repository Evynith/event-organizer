package main

import (
	"main/internal/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	public := router.Group("/events")
	public.Use()
	{
		public.GET("", controller.Events)   //status 0 only admin
		public.GET(":id", controller.Event) //status 0 only admin
		private := public.Group("")
		private.Use()
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

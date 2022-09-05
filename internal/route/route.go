package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	controller "main/internal/controller"
	"main/internal/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to API event-organizer"})
	})

	router.POST("/login", func(ctx *gin.Context) {
		token := controller.Login(ctx)

		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to create access"})
		}
	})

	var middle middleware.JWTmiddleware = middleware.JWTServiceMiddleware("admin")
	var eController controller.EventControllerInterface = controller.EventControllerStart(middle)
	var iController controller.InscriptionControllerInterface = controller.InscriptionControllerStart(middle)

	public := router.Group("/events")
	public.Use(middle.AuthorizeJWT())
	{
		public.GET("", eController.Events)
		public.GET(":id", eController.Event)
		private := public.Group("")
		private.Use(middle.OnlyUser())
		{
			private.POST("", eController.PostEvent)
			private.DELETE(":id", eController.DeleteEvent)
			private.PUT(":id", eController.PutEvent)
		}

		inscription := public.Group("inscription")
		inscription.Use()
		{
			inscription.POST("", iController.PostInscription)
			inscription.GET("", iController.Inscriptions)
		}
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})
	return router
}

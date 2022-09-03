package route

import (
	controller "main/internal/controller"
	"main/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/login", func(ctx *gin.Context) {
		token := controller.Login(ctx)

		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, "No se pudo crear acceso")
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
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	return router
}

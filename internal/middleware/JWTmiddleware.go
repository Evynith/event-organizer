package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	service "main/internal/service"
)

type JWTmiddleware interface {
	AuthorizeJWT() gin.HandlerFunc
	OnlyUser() gin.HandlerFunc
	GetType() string
}
type jwtServices struct {
	espectedUser string
	typeUser     string
	token        string
	id           string
}

func (middleware *jwtServices) GetType() string {
	return middleware.typeUser
}

func JWTServiceMiddleware(espected string) JWTmiddleware {
	return &jwtServices{
		espectedUser: espected,
		token:        "",
		typeUser:     "",
		id:           "",
	}
}

/*
Obtiene de la cabecera de la petición el token recibido y a través de él el tipo de usuario y su id.
Revisa que el token sea válido y, de serlo, que haya sido guardado con anterioridad en la base de datos
*/
func (middleware *jwtServices) AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := ctx.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		middleware.token = tokenString
		token, err := service.JWTAuthService().ValidateToken(middleware.token)

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		middleware.typeUser, middleware.id = service.JWTAuthService().TypeUser(token)
		if middleware.typeUser != "" && middleware.id != "" {
			returnAs(ctx, (service.ExistsToken(middleware.token, middleware.id)))
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}

/*
Revisa que la ruta que lo utiliza sea accedida sólo por el tipo de usuario indicado
al inicializar el elemento de restricción
*/
func (middleware *jwtServices) OnlyUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		returnAs(ctx, (middleware.espectedUser == middleware.typeUser))
	}
}

/*
Según el valor recibido por una condición elige si continuar ó abortar la función de restricción de acceso padre
*/
func returnAs(ctx *gin.Context, status bool) {
	if status {
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

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

func (middleware *jwtServices) AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := ctx.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		middleware.token = tokenString
		token, err := service.JWTAuthService().ValidateToken(middleware.token)
		middleware.typeUser, middleware.id = service.JWTAuthService().TypeUser(middleware.token)

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			returnAs(ctx, (service.ExistsToken(middleware.token)))
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}

func (middleware *jwtServices) OnlyUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		returnAs(ctx, (middleware.espectedUser == middleware.typeUser))
	}
}

func returnAs(ctx *gin.Context, status bool) {
	if status {
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

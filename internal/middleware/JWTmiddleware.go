package middleware

import (
	service "main/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTmiddleware interface {
	AuthorizeJWT() gin.HandlerFunc
	OnlyUser() gin.HandlerFunc
}
type jwtServices struct {
	typeUser string
	token    string
}

func JWTServiceMiddleware(typeU string) JWTmiddleware {
	return &jwtServices{
		typeUser: typeU,
	}
}

func (middleware *jwtServices) AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		middleware.token = tokenString
		token, _ := service.JWTAuthService().ValidateToken(middleware.token)
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if service.ExistsToken(middleware.token) {
				c.Next()
			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}

func (middleware *jwtServices) OnlyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if service.JWTAuthService().ValidateUser(middleware.token, middleware.typeUser) == true {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
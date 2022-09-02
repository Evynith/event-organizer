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
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		middleware.token = tokenString
		token, err := service.JWTAuthService().ValidateToken(middleware.token)
		middleware.typeUser, middleware.id = service.JWTAuthService().TypeUser(middleware.token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
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
		if middleware.espectedUser == middleware.typeUser {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

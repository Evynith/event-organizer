package middleware

import (
	"fmt"
	"log"

	service "main/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		//log.Println(tokenString)
		token, err := service.JWTAuthService().ValidateToken(tokenString)
		log.Println(token.Claims.(jwt.MapClaims))
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if service.ExistsToken(tokenString) {
				fmt.Println(claims["name"])
				c.Next()
			}
		}
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func OnlyAdmin(user string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		if service.JWTAuthService().ValidateUser(tokenString, "admin") == true {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

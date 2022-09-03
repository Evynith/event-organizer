package controller

import (
	service "main/internal/service"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) string {
	var jwtService service.JWTService = service.JWTAuthService()
	var token string = ""
	user, password, hasAuth := ctx.Request.BasicAuth()
	if !hasAuth {
		return ""
	}
	isUserAuthenticated := service.LoginUser(user, password)
	if isUserAuthenticated {
		token = jwtService.GenerateToken(user)
	}
	service.SaveToken(user, token)

	return token
}

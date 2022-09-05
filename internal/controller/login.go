package controller

import (
	"github.com/gin-gonic/gin"

	service "main/internal/service"
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

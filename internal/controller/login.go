package controller

import (
	"main/internal/model"
	service "main/internal/service"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) string {
	var jwtService service.JWTService = service.JWTAuthService()
	var credential model.User
	var token string = ""
	err := ctx.ShouldBind(&credential)
	if err != nil {
		return "no data found"
	}
	isUserAuthenticated := service.LoginUser(credential.Username, credential.Password)
	if isUserAuthenticated {
		token = jwtService.GenerateToken(credential.Username)
	}
	service.SaveToken(credential.Username, token)

	return token
}

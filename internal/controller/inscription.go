package controller

import (
	"main/internal/model"
	inscriptionRepository "main/internal/repository/inscription"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostInscription(c *gin.Context) {
	var newInscription model.Inscription
	if err := c.BindJSON(&newInscription); err != nil {
		return
	}

	err := inscriptionRepository.Create(newInscription)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, newInscription)
}

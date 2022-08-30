package controller

import (
	"errors"
	"main/internal/model"
	eventRepository "main/internal/repository/event"
	inscriptionRepository "main/internal/repository/inscription"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostInscription(c *gin.Context) {
	var newInscription model.Inscription
	err := c.BindJSON(&newInscription)
	idEvent := newInscription.Event.Hex()
	event, err := eventRepository.ReadOne(idEvent)

	if event.DateOfEvent < primitive.NewDateTimeFromTime(time.Now()) {
		c.AbortWithError(http.StatusNotAcceptable, errors.New("This event as been finalized"))
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = inscriptionRepository.Create(newInscription)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, newInscription)
}

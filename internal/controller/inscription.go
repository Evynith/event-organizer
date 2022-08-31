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

func Inscriptions(c *gin.Context) {
	title, _ := c.GetQuery("title")
	date1, _ := c.GetQuery("since")
	date2, _ := c.GetQuery("until")
	state, _ := c.GetQuery("state")

	var newInscription model.Inscription
	err := c.BindJSON(&newInscription)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	inscriptions, err := inscriptionRepository.Read(newInscription.User)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var inscriptionsList []primitive.ObjectID
	for _, inscriptionID := range inscriptions {
		inscriptionsList = append(inscriptionsList, inscriptionID.Event)
	}

	events, err := eventRepository.Read(title, date1, date2, state, inscriptionsList)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, events)
}

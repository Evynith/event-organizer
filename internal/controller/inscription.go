package controller

import (
	"errors"
	mi "main/internal/middleware"
	"main/internal/model"
	eventRepository "main/internal/repository/event"
	inscriptionRepository "main/internal/repository/inscription"
	"main/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InscriptionControllerInterface interface {
	PostInscription(c *gin.Context)
	Inscriptions(c *gin.Context)
}

type inscriptionController struct {
	middle mi.JWTmiddleware
}

func InscriptionControllerStart(j mi.JWTmiddleware) InscriptionControllerInterface {
	return &inscriptionController{
		middle: j,
	}
}

func (i *inscriptionController) PostInscription(c *gin.Context) {
	var newInscription model.Inscription
	err := c.BindJSON(&newInscription)
	idEvent := newInscription.Event.Hex()
	access := service.AccessDraft(i.middle.GetType())
	idOld := service.CreateFilterEvent(idEvent, access)
	event, err := eventRepository.ReadOne(idOld)

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

func (i *inscriptionController) Inscriptions(c *gin.Context) {
	var newFilter model.Filter
	if err := c.BindQuery(&newFilter); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	access := service.AccessDraft(i.middle.GetType())
	filter := service.CreateFilterEvents(newFilter, []primitive.ObjectID{}, access)

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

	events, err := eventRepository.Read(filter)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, events)
}

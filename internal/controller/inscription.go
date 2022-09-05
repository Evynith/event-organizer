package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	middleware "main/internal/middleware"
	"main/internal/model"
	eventRepository "main/internal/repository/event"
	inscriptionRepository "main/internal/repository/inscription"
	"main/internal/service"
)

type InscriptionControllerInterface interface {
	PostInscription(ctx *gin.Context)
	Inscriptions(ctx *gin.Context)
}

type inscriptionController struct {
	handler middleware.JWTmiddleware
}

func InscriptionControllerStart(middle middleware.JWTmiddleware) InscriptionControllerInterface {
	return &inscriptionController{
		handler: middle,
	}
}

func (i *inscriptionController) PostInscription(ctx *gin.Context) {
	var newInscription model.Inscription
	bindJSON(ctx, &newInscription)
	idEvent := newInscription.Event.Hex()
	access := service.AccessDraft(i.handler.GetType())
	idOld := service.CreateFilterEvent(idEvent, access)

	event, err := eventRepository.ReadOne(idOld)
	inspectError(ctx, err)

	if event.DateOfEvent < primitive.NewDateTimeFromTime(time.Now()) {
		ctx.AbortWithError(http.StatusNotAcceptable, errors.New("This event as been finalized"))
		return
	}

	err = inscriptionRepository.Create(newInscription)
	inspectError(ctx, err)

	ctx.IndentedJSON(http.StatusCreated, newInscription)
}

func (i *inscriptionController) Inscriptions(ctx *gin.Context) {
	var newFilter model.Filter
	var newInscription model.Inscription
	var inscriptionsList []primitive.ObjectID
	bindQuery(ctx, &newFilter)
	bindJSON(ctx, &newInscription)

	access := service.AccessDraft(i.handler.GetType())
	filter := service.CreateFilterEvents(newFilter, []primitive.ObjectID{}, access)

	inscriptions, err := inscriptionRepository.Read(newInscription.User)
	inspectError(ctx, err)

	for _, inscriptionID := range inscriptions {
		inscriptionsList = append(inscriptionsList, inscriptionID.Event)
	}

	filter = service.CreateFilterListOfEvent(filter, inscriptionsList)
	events, err := eventRepository.Read(filter)
	inspectError(ctx, err)

	ctx.IndentedJSON(http.StatusOK, events)
}

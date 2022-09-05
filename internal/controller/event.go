package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	middleware "main/internal/middleware"
	"main/internal/model"
	eventRepository "main/internal/repository/event"
	"main/internal/service"
)

type EventControllerInterface interface {
	Events(ctx *gin.Context)
	Event(ctx *gin.Context)
	PostEvent(ctx *gin.Context)
	DeleteEvent(ctx *gin.Context)
	PutEvent(ctx *gin.Context)
}

type eventController struct {
	handler middleware.JWTmiddleware
}

func EventControllerStart(middle middleware.JWTmiddleware) EventControllerInterface {
	return &eventController{
		handler: middle,
	}
}

func (e *eventController) Events(ctx *gin.Context) {
	var newFilter model.Filter
	bindQuery(ctx, &newFilter)
	access := service.AccessDraft(e.handler.GetType())
	filter := service.CreateFilterEvents(newFilter, []primitive.ObjectID{}, access)

	usuarios, err := eventRepository.Read(filter)
	inspectError(ctx, err)

	ctx.IndentedJSON(http.StatusOK, usuarios)
}

func (e *eventController) Event(ctx *gin.Context) {
	id := ctx.Param("id")
	access := service.AccessDraft(e.handler.GetType())
	idOld := service.CreateFilterEvent(id, access)

	user, err := eventRepository.ReadOne(idOld)
	inspectError(ctx, err)

	ctx.IndentedJSON(http.StatusOK, user)
}

func (e *eventController) PostEvent(ctx *gin.Context) {
	var newEvent model.Event
	bindJSON(ctx, &newEvent)

	result, err := eventRepository.Create(newEvent)
	inspectError(ctx, err)
	newEvent.ID = result.(primitive.ObjectID)

	ctx.IndentedJSON(http.StatusCreated, newEvent)
}

func (e *eventController) DeleteEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	idOld := service.CreateFilterID(id)

	err := eventRepository.Delete(idOld)
	inspectError(ctx, err)

	ctx.Status(200)
}

func (e *eventController) PutEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	var newEvent model.Event
	bindJSON(ctx, &newEvent)
	data := service.CreateEventUpdate(newEvent)
	idOld := service.CreateFilterID(id)

	err := eventRepository.Update(data, idOld)
	inspectError(ctx, err)

	ctx.IndentedJSON(http.StatusOK, newEvent)
}

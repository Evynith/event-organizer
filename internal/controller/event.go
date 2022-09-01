package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mi "main/internal/middleware"
	"main/internal/model"
	eventRepository "main/internal/repository/event"
	"main/internal/service"
)

type EventControllerInterface interface {
	Events(c *gin.Context)
	Event(c *gin.Context)
	PostEvent(c *gin.Context)
	DeleteEvent(c *gin.Context)
	PutEvent(c *gin.Context)
}

type eventController struct {
	middle mi.JWTmiddleware
}

func EventControllerStart(j mi.JWTmiddleware) EventControllerInterface {
	return &eventController{
		middle: j,
	}
}

func (e *eventController) Events(c *gin.Context) {
	var newFilter model.Filter
	if err := c.BindQuery(&newFilter); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	access := service.AccessDraft(e.middle.GetType())
	filter := service.CreateFilterEvents(newFilter, []primitive.ObjectID{}, access)

	usuarios, err := eventRepository.Read(filter)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, usuarios)
}

func (e *eventController) Event(c *gin.Context) {
	id := c.Param("id")
	access := service.AccessDraft(e.middle.GetType())
	idOld := service.CreateFilterEvent(id, access)
	user, err := eventRepository.ReadOne(idOld)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func (e *eventController) PostEvent(c *gin.Context) {
	var newEvent model.Event
	if err := c.BindJSON(&newEvent); err != nil {
		return
	}

	result, err := eventRepository.Create(newEvent)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	newEvent.ID = result.(primitive.ObjectID)
	c.IndentedJSON(http.StatusCreated, newEvent)
}

func (e *eventController) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	idOld := service.CreateFilterID(id)
	err := eventRepository.Delete(idOld)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(200)
}

func (e *eventController) PutEvent(c *gin.Context) {
	id := c.Param("id")
	var newEvent model.Event
	if err := c.BindJSON(&newEvent); err != nil {
		return
	}
	data := service.CreateEventUpdate(newEvent)
	idOld := service.CreateFilterID(id)
	err := eventRepository.Update(data, idOld)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, newEvent)
}

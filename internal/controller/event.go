package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"main/internal/model"
	eventRepository "main/internal/repository/event"
)

func Events(c *gin.Context) {
	title, _ := c.GetQuery("title")
	date1, _ := c.GetQuery("since")
	date2, _ := c.GetQuery("until")
	state, _ := c.GetQuery("state")

	usuarios, err := eventRepository.Read(title, date1, date2, state, []primitive.ObjectID{})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, usuarios)
}

func Event(c *gin.Context) {
	id := c.Param("id")
	user, err := eventRepository.ReadOne(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func PostEvent(c *gin.Context) {
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

func DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	err := eventRepository.Delete(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(200)
}

func PutEvent(c *gin.Context) {
	id := c.Param("id")
	var newEvent model.Event
	if err := c.BindJSON(&newEvent); err != nil {
		return
	}
	err := eventRepository.Update(newEvent, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, newEvent)
}

package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/surajgoraicse/go_event/internal/database"
)

func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := app.models.Events.Insert(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event", "err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)

}



func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"invalid id passed": c.Param("id")})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	updateEvent := &database.Event{}
	if err := c.ShouldBindJSON(updateEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateEvent.ID = id

	if err := app.models.Events.Update(updateEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updateEvent)

}

func (app *application) deletedEvent(c *gin.Context){
	id , err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid id passed"})
	}

	if err :=app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "failed to delete"})
		return
	}

	c.JSON(http.StatusNoContent, nil)

}
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

	// binding the context user (from authorization) as the owner_id
	// authorization : 
	contextUser := app.GetUserFromContext(c)
	event.OwnerID = contextUser.ID 


	if err := app.models.Events.Insert(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event", "err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)

}

// getAllEvents returns all events
//
// @Summary Returns all events
// @Description Returns all events
// @Tags Events
// @Accept json
// @Produce json
// @Success 200  {object} []database.Event
// @Router /api/v1/events [get]

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

	user := app.GetUserFromContext(c)

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	if event.OwnerID != user.ID{
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update the event"})
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

func (app *application) deletedEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id passed"})
		return
	}

	user := app.GetUserFromContext(c)
	event , err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching event using id"})
		return 

	}
	if event == nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "event not found"})
		return
	}

	
	if user.ID != event.OwnerID{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not authorized to delete the event"})
		return 
	}
	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	event, err := app.models.Events.Get(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// authorization : 
	contextUser := app.GetUserFromContext(c)
	if contextUser.ID != event.OwnerID{
		c.JSON(http.StatusUnauthorized, gin.H{"error":"user is not authorized "})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err, "message": "event not found"})
		return
	}
	user, err := app.models.Users.Get(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err, "message": "user not found"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(eventID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "attendee already exists"})
		return
	}
	attendee := database.Attendees{
		UserID:  userID,
		EventID: eventID,
	}
	if _, err = app.models.Attendees.Insert(&attendee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err, "message": "error inserting into db"})
		return
	}
	c.JSON(http.StatusCreated, attendee)

}

func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err, "message": "invalid event id passed"})
		return
	}

	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err, "message": "failed to retrieve"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	eventID, err := strconv.Atoi(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	event , err := app.models.Events.Get(eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	// authorization : 
	contextUser := app.GetUserFromContext(c)
	if contextUser.ID != event.OwnerID{
		c.JSON(http.StatusUnauthorized, gin.H{"error":"user is not authorized "})
		return
	}

	err = app.models.Attendees.Delete(userID, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

func (app *application) getEventsByAttendee(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	events, err := app.models.Attendees.GetEventsByAttendee(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get event"})
	}
	c.JSON(http.StatusOK, events)

}

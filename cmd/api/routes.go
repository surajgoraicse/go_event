package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")

	{
		v1.POST("/events", app.createEvent)
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)
		v1.PUT("/events/:id", app.updateEvent)
		v1.DELETE("/events", app.deletedEvent)

		v1.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
		v1.GET("/events/:id/attendees", app.getAttendeesForEvent)
		v1.DELETE("/events/:id/atteendees/:userID", app.deleteAttendeeFromEvent)

		v1.GET("/attendees/:id/events", app.getEventsByAttendee) // returns all events participating by an attendee
		v1.POST("/auth/register", app.registerUser)
	}

	return g

}

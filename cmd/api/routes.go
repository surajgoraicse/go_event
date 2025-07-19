package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")

	{
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)

		v1.GET("/events/:id/attendees", app.getAttendeesForEvent)

		v1.GET("/attendees/:id/events", app.getEventsByAttendee) // returns all events participating by an attendee
		v1.POST("/auth/register", app.registerUser)

		v1.POST("/auth/login", app.login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		v1.POST("/events", app.createEvent)
		v1.PUT("/events/:id", app.updateEvent)
		v1.DELETE("/events", app.deletedEvent)
		v1.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
		v1.DELETE("/events/:id/atteendees/:userID", app.deleteAttendeeFromEvent)

	}
	g.GET("/", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{"success" : "true"})
	})

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(302, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8000/swagger/doc.json"))(c)

	})

	return g

}

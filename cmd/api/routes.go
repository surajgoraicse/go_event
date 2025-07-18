package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")

	{
		v1.POST("/evets", app.createEvent)
		v1.GET("/evets", app.getAllEvents)
		v1.GET("/evets/:id", app.getEvent)
		v1.PUT("/evets/:id", app.updateEvent)
		v1.DELETE("/evets", app.deletedEvent)
	}

	return g

}


func (m *)
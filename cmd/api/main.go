package main

import (
	"database/sql"
	"log"

	_ "github.com/surajgoraicse/go_event/docs"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/surajgoraicse/go_event/internal/database"
	"github.com/surajgoraicse/go_event/internal/env"
)

// @title Go Gin Resp API
// @version 1.0
// @description A rest API in GO using Gin framework
// @in header
// @name header
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**


type application struct {
	port      int
	jwtSecret string
	models    *database.Models
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123234"),
		models:    models,
	}
	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

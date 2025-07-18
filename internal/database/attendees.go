package database

import "database/sql"

type AttendeesModel struct {
	DB *sql.DB
}

type Attendees struct {
	ID      int `json:"id"`
	UserID  int `json:"userID"`
	EventID int `json:"eventID"`
}

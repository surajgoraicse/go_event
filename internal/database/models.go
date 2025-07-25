package database

import "database/sql"

type Models struct {
	Users     UserModel
	Events    EventModel
	Attendees AttendeesModel
}
// constructor function 
func NewModels(db *sql.DB) *Models {
	return &Models{
		Users:     UserModel{DB: db},
		Events:    EventModel{DB: db},
		Attendees: AttendeesModel{DB: db},
	}
}

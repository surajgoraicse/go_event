package database

import (
	"context"
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	ID          int       `json:"id"`
	OwnerID     int       `json:"ownerId" binding:"required"` // used by gin
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description" binding:"required,min=10"`
	Date        time.Time `json:"date" binding:"required, datetime=2006-01-02"`
	Location    string    `json:"location" binding:"required, min=3"`
}

// NOTE: what are binding tags

func (m *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // NOTE: what is context here and it uses
	defer cancel()

	query := "INSERT INTO events (owner_id, name, description, date, location) VALUES ($1, $2, $3, $4,$5)"

	// what is happening here
	row := m.DB.QueryRowContext(ctx, query, event.OwnerID, event.Name, event.Description, event.Date, event.Location)

	return row.Scan(&event.ID)
}

func (m *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM events"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.OwnerID, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil

}

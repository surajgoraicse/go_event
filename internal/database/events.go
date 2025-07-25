package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	ID          int       `json:"id"`
	OwnerID     int       `json:"ownerId" binding:"required"` // (made changes here : remove required)
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description" binding:"required,min=10"`
	Date        time.Time `json:"date" binding:"required"`
	Location    string    `json:"location" binding:"required,min=3"`
}

// NOTE: what are binding tags

func (m *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // NOTE: what is context here and it uses
	defer cancel()

	query := `
		INSERT INTO events (owner_id, name, description, date, location)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id` // ✅ Required

	// what is happening here

	// return m.DB.QueryRowContext(ctx, query,
	// 	event.OwnerID,
	// 	event.Name,
	// 	event.Description,
	// 	event.Date,
	// 	event.Location,
	// ).Scan(&event.ID)
	// NOTE: what is happening here

	err := m.DB.QueryRowContext(ctx, query,
		event.OwnerID,
		event.Name,
		event.Description,
		event.Date,
		event.Location,
	).Scan(&event.ID)
	// NOTE: what is happening here
	if err != nil {
		fmt.Println(err)
	}
	return err
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

func (m *EventModel) Get(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM events WHERE ID = $1"
	// var event *Event
	event := new(Event)

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&event.ID, &event.OwnerID, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}

	}
	return event, nil
}

func (m *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE events SET name = $1, description = $2, date = $3, location = $4 WHERE id = $5"

	_, err := m.DB.ExecContext(ctx, query, event.Name, event.Description, event.Date, event.Location, event.ID)

	if err != nil {
		return err
	}
	return nil
}
func (m *EventModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE events where id = $1"
	_, err := m.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

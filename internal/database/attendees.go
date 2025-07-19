package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeesModel struct {
	DB *sql.DB
}

type Attendees struct {
	ID      int `json:"id"`
	UserID  int `json:"userID"`
	EventID int `json:"eventID"`
}

func (m *AttendeesModel) Insert(attendee *Attendees) (*Attendees, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO attendees (event_id, user_id) VALUES ($1, $2) RETURNING id"
	if err := m.DB.QueryRowContext(ctx, query, attendee.EventID, attendee.UserID).Scan(&attendee.ID); err != nil {
		return nil, err
	}
	return attendee, nil
}

func (m *AttendeesModel) GetByEventAndAttendee(eventID, userID int) (*Attendees, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM attendees WHERE event_id = $1 and user_id = $2"

	// var attendee *Attendees // this will give error as it is just declaring the variable but not allocating the memory
	// attendee := new(Attendees) // replacement 1
	attendee := &Attendees{} // replacement 2
	if err := m.DB.QueryRowContext(ctx, query, eventID, userID).Scan(&attendee.ID, &attendee.UserID, &attendee.EventID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return attendee, nil
}

func (m *AttendeesModel) GetAttendeesByEvent(id int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT u.id, u.name, u.email
		FROM users u
		JOIN attendees a ON u.id = a.user_id
		where a.event_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// delete attendee from event
func (m *AttendeesModel) Delete(userID, eventID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM attendees WHERE user_id = $1 and event_id = $2"
	return m.DB.QueryRowContext(ctx, query, userID, eventID).Scan()
}

// list all events of the user (attending)
func (m *AttendeesModel) GetEventsByAttendee(userID int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.user_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	events := []*Event{}
	defer rows.Close()
	for rows.Next() {
		event := &Event{}
		if err :=rows.Scan(&event.ID, &event.OwnerID, &event.Name, &event.Description, &event.Date, &event.Location); err != nil {
			return nil , err
		}
		events = append(events, event)
	}
	return events, nil
}

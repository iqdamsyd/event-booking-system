package repository

import (
	"context"
	"event-booking-system/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) List(ctx context.Context) ([]models.Event, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, description, location, event_date, capacity, booked_count, created_at, updated_at
		FROM events
		ORDER BY event_date ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.EventDate, &event.Capacity, &event.BookedCount, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetByID(ctx context.Context, id string) (*models.Event, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, description, location, event_date, capacity, booked_count, created_at, updated_at
		FROM events
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var event models.Event
	for rows.Next() {
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.EventDate, &event.Capacity, &event.BookedCount, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, err
		}
	}

	return &event, nil
}

func (r *EventRepository) Create(ctx context.Context, event models.Event) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO events (title, description, location, event_date, capacity)
		VALUES ($1, $2, $3, $4, $5)
	`, event.Title, event.Description, event.Location, event.EventDate, event.Capacity)

	return err
}

func (r *EventRepository) Update(ctx context.Context, event models.Event) error {
	_, err := r.db.Exec(ctx, `
		UPDATE events SET title = $1, description = $2, location = $3, event_date = $4, capacity = $5, updated_at = NOW()
		WHERE id = $6
	`, event.Title, event.Description, event.Location, event.EventDate, event.Capacity, event.ID)

	return err
}

func (r *EventRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM events WHERE id = $1
	`, id)

	return err
}

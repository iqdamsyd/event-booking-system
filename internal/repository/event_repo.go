package repository

import (
	"context"
	"event-booking-system/internal/helper"
	"event-booking-system/internal/models"
	"fmt"

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

func (r *EventRepository) List(ctx context.Context, filter models.RequestFilterEvent) ([]models.Event, *models.Meta, error) {
	query := `
		SELECT id, title, description, location, event_date, capacity, booked_count, created_at, updated_at
		FROM events
		WHERE is_deleted IS NULL
	`
	var params []interface{}

	if filter.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", len(params)+1)
		params = append(params, "%"+fmt.Sprintf("%s", filter.Title)+"%")
	}

	if filter.Location != "" {
		query += fmt.Sprintf(" AND location ILIKE $%d", len(params)+1)
		params = append(params, "%"+fmt.Sprintf("%s", filter.Location)+"%")
	}

	if filter.StartDate != "" && filter.EndDate != "" {
		parsedStartDate, err := helper.ParseYYYYMMDD(filter.StartDate)
		if err != nil {
			return nil, nil, err
		}
		query += fmt.Sprintf(" AND event_date > $%d", len(params)+1)
		params = append(params, parsedStartDate)

		parsedEndDate, err := helper.ParseYYYYMMDD2359(filter.EndDate)
		if err != nil {
			return nil, nil, err
		}
		query += fmt.Sprintf(" AND event_date < $%d", len(params)+1)
		params = append(params, parsedEndDate)
	}

	if filter.AvailableOnly {
		query += ` AND booked_count < capacity`
	}

	if filter.SortBy != "" && filter.SortType != "" {
		query += fmt.Sprintf(" ORDER BY $%d %s", len(params)+1, filter.SortType)
		params = append(params, filter.SortBy)
	}

	if filter.Page != 0 && filter.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(params)+1)
		params = append(params, filter.Limit)
		query += fmt.Sprintf(" OFFSET $%d", len(params)+1)
		params = append(params, (filter.Page*filter.Limit)-filter.Limit)
	}

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Location, &event.EventDate, &event.Capacity, &event.BookedCount, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, nil, err
		}
		events = append(events, event)
	}

	countQuery := "SELECT COUNT(*) FROM events WHERE is_deleted IS NULL"
	var countParams []interface{}

	if filter.Title != "" {
		countQuery += fmt.Sprintf(" AND title ILIKE $%d", len(countParams)+1)
		countParams = append(countParams, "%"+fmt.Sprintf("%s", filter.Title)+"%")
	}

	if filter.Location != "" {
		countQuery += fmt.Sprintf(" AND location ILIKE $%d", len(countParams)+1)
		countParams = append(countParams, "%"+fmt.Sprintf("%s", filter.Location)+"%")
	}

	if filter.StartDate != "" && filter.EndDate != "" {
		parsedStartDate, err := helper.ParseYYYYMMDD(filter.StartDate)
		if err != nil {
			return nil, nil, err
		}
		countQuery += fmt.Sprintf(" AND event_date > $%d", len(countParams)+1)
		countParams = append(countParams, parsedStartDate)

		parsedEndDate, err := helper.ParseYYYYMMDD2359(filter.EndDate)
		if err != nil {
			return nil, nil, err
		}
		countQuery += fmt.Sprintf(" AND event_date < $%d", len(countParams)+1)
		countParams = append(countParams, parsedEndDate)
	}

	if filter.AvailableOnly {
		countQuery += ` AND booked_count < capacity`
	}

	var count int
	row := r.db.QueryRow(ctx, countQuery, countParams...)
	if err := row.Scan(&count); err != nil {
		return nil, nil, err
	}

	var meta models.Meta
	meta.CurrentPage = filter.Page
	meta.PageSize = filter.Limit
	meta.TotalData = count
	meta.LastPage = count / filter.Limit
	if count%filter.Limit != 0 {
		meta.LastPage++
	}
	if meta.LastPage < meta.CurrentPage {
		meta.CurrentPage = meta.LastPage
	}

	return events, &meta, nil
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

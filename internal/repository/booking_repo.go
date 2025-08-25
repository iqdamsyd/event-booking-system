package repository

import (
	"context"
	"event-booking-system/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (r *BookingRepository) GetAllByUserID(ctx context.Context, userID string) ([]models.MyBooking, error) {
	query := `
		SELECT b.id, e.title, e.location, e.event_date, b.seats, b.status, b.created_at, b.updated_at
		FROM bookings b
		LEFT JOIN events e on e.id = b.event_id
		WHERE b.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var bookings []models.MyBooking
	for rows.Next() {
		var booking models.MyBooking
		if err := rows.Scan(&booking.ID, &booking.Title, &booking.Location, &booking.EventDate, &booking.Seats, &booking.Status, &booking.CreatedAt, &booking.UpdatedAt); err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, event_id, seats, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var booking models.Booking
	for rows.Next() {
		if err := rows.Scan(&booking.ID, &booking.UserID, &booking.EventID, &booking.Seats, &booking.Status, &booking.CreatedAt, &booking.UpdatedAt); err != nil {
			return nil, err
		}
	}

	return &booking, nil
}

func (r *BookingRepository) Create(ctx context.Context, booking models.Booking) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	exists, _ := tx.Query(ctx, `
		SELECT id, user_id, event_id, seats, status, created_at, updated_at
		FROM bookings
		WHERE event_id = $1 AND seats = $2 AND status != 'cancelled'::booking_statuses
	`, booking.EventID, booking.Seats)
	defer exists.Close()
	if exists.Next() {
		return fmt.Errorf("get booking: seat is booked")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO bookings (user_id, event_id, seats)
		VALUES ($1, $2, $3)
	`, booking.UserID, booking.EventID, booking.Seats)
	if err != nil {
		return fmt.Errorf("insert booking: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE events
		SET booked_count = booked_count + 1, updated_at = NOW()
		WHERE id = $1
	`, booking.EventID)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *BookingRepository) Cancel(ctx context.Context, id string, eventID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE bookings
		SET status = 'cancelled'::booking_statuses, updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("update booking: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE events
		SET booked_count = booked_count - 1, updated_at = NOW()
		WHERE id = $1
	`, eventID)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *BookingRepository) Confirm(ctx context.Context, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE bookings
		SET status = 'confirmed'::booking_statuses, updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("update booking: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *BookingRepository) CountBookingStatusByEventID(ctx context.Context, eventID string) ([]models.CountBookingStatus, error) {
	query := `
		SELECT event_id, status, COUNT(*) AS count
		FROM bookings
		GROUP BY event_id, status
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var bookings []models.CountBookingStatus
	for rows.Next() {
		var booking models.CountBookingStatus
		if err := rows.Scan(&booking.EventID, &booking.Status, &booking.Count); err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

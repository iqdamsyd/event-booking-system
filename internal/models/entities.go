package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"` // hashed
	Role      string    `db:"role"`     // e.g. "user", "admin"
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Event struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Location    string    `db:"location"`
	EventDate   time.Time `db:"event_date"`
	Capacity    int       `db:"capacity"`
	BookedCount int       `db:"booked_count"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Booking struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	EventID   uuid.UUID `db:"event_id"`
	Seats     int       `db:"seats"`
	Status    string    `db:"status"` // e.g. "pending", "confirmed", "cancelled"
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

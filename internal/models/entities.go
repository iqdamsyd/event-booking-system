package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"` // hashed
	Role      string    `db:"role" json:"role"`         // e.g. "user", "admin"
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Event struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Location    string    `db:"location" json:"location"`
	EventDate   time.Time `db:"event_date" json:"event_date"`
	Capacity    int       `db:"capacity" json:"capacity"`
	BookedCount int       `db:"booked_count" json:"booked_count"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type Booking struct {
	ID        uuid.UUID `db:"id" json:"id" `
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	EventID   uuid.UUID `db:"event_id" json:"event_id"`
	Seats     int       `db:"seats" json:"seats"`
	Status    string    `db:"status" json:"status"` // e.g. "pending", "confirmed", "cancelled"
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

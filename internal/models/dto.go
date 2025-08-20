package models

import (
	"time"

	"github.com/google/uuid"
)

/* User */
type RequestCreateUser struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type RequestUpdateUser struct {
	ID       uuid.UUID `param:"id"`
	Name     string    `form:"name"`
	Email    string    `form:"email"`
	Password string    `form:"password"`
}

type RequestRegisterUser struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type RequestLoginUser struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

/* Event */
type RequestCreateEvent struct {
	Title       string    `form:"title"`
	Description string    `form:"description"`
	Location    string    `form:"location"`
	EventDate   time.Time `form:"event_date"`
	Capacity    int       `form:"capacity"`
}

type RequestUpdateEvent struct {
	ID          uuid.UUID `param:"id"`
	Title       string    `form:"title"`
	Description string    `form:"description"`
	Location    string    `form:"location"`
	EventDate   time.Time `form:"event_date"`
	Capacity    int       `form:"capacity"`
}

/* Booking */
type RequestBooking struct {
	UserID  uuid.UUID
	EventID uuid.UUID `json:"event_id"`
	Seats   int       `json:"seats"`
}

type RequestCancelBooking struct {
	UserID uuid.UUID
	ID     uuid.UUID `json:"id"`
}

type RequestConfirmBooking struct {
	UserID uuid.UUID
	ID     uuid.UUID `json:"id"`
}

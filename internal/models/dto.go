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

type RequestFilterEvent struct {
	Page          int    `query:"page"`
	Limit         int    `query:"limit"`
	SortBy        string `query:"sort_by"`
	SortType      string `query:"sort_type"`
	Title         string `query:"title"`
	Location      string `query:"location"`
	StartDate     string `query:"start_date"`
	EndDate       string `query:"end_date"`
	AvailableOnly bool   `query:"available_only"`
}

func NewRequestFilterEvent() RequestFilterEvent {
	return RequestFilterEvent{
		Page:          1,
		Limit:         10,
		SortBy:        "created_at",
		SortType:      "desc",
		AvailableOnly: false,
	}
}

type EventOverview struct {
	Event
	Count struct {
		Pending   int `json:"pending"`
		Confirmed int `json:"confirmed"`
		Cancelled int `json:"cancelled"`
	} `json:"count"`
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

type MyBooking struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Location  string    `json:"location"`
	EventDate time.Time `json:"event_date"`
	Seats     int       `json:"seats"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CountBookingStatus struct {
	EventID uuid.UUID `json:"event_id"`
	Status  string    `json:"status"`
	Count   int       `json:"count"`
}

/* Meta */
type Meta struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalData   int `json:"total_data"`
	LastPage    int `json:"last_page"`
}

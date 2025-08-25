package service

import (
	"context"
	"errors"
	"event-booking-system/internal/models"
	"event-booking-system/internal/repository"
	"fmt"

	"github.com/google/uuid"
)

const (
	PENDING   = "pending"
	CANCELLED = "cancelled"
	CONFIRMED = "confirmed"
)

type BookingService struct {
	repo      *repository.BookingRepository
	eventRepo *repository.EventRepository
}

func NewBookingService(repo *repository.BookingRepository, eventRepo *repository.EventRepository) *BookingService {
	return &BookingService{
		repo:      repo,
		eventRepo: eventRepo,
	}
}

func (s *BookingService) Create(ctx context.Context, payload models.RequestBooking) error {
	event, err := s.eventRepo.GetByID(ctx, payload.EventID.String())
	if err != nil {
		return err
	}

	if event.BookedCount >= event.Capacity {
		return errors.New("event is full")
	}

	var booking models.Booking
	booking.EventID = payload.EventID
	booking.UserID = payload.UserID
	booking.Seats = payload.Seats
	if err := s.repo.Create(ctx, booking); err != nil {
		return err
	}

	return nil
}

func (s *BookingService) Cancel(ctx context.Context, payload models.RequestCancelBooking) error {
	booking, err := s.repo.GetByID(ctx, payload.ID.String())
	if err != nil {
		return err
	}

	if booking.UserID != payload.UserID {
		return fmt.Errorf("unauthorized")
	}

	if booking.Status != PENDING {
		return fmt.Errorf("can not cancel")
	}

	if err := s.repo.Cancel(ctx, payload.ID.String(), booking.EventID.String()); err != nil {
		return err
	}

	return nil
}

func (s *BookingService) Confim(ctx context.Context, payload models.RequestConfirmBooking) error {
	booking, err := s.repo.GetByID(ctx, payload.ID.String())
	if err != nil {
		return err
	}

	if booking.UserID != payload.UserID {
		return fmt.Errorf("unauthorized")
	}

	if booking.Status != PENDING {
		return fmt.Errorf("can not cancel")
	}

	if err := s.repo.Confirm(ctx, payload.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *BookingService) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.MyBooking, error) {
	bookings, err := s.repo.GetAllByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

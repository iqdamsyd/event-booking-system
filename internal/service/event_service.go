package service

import (
	"context"
	"event-booking-system/internal/models"
	"event-booking-system/internal/repository"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Create(ctx context.Context, payload models.RequestCreateEvent) error {
	var event models.Event
	event.Title = payload.Title
	event.Description = payload.Description
	event.Location = payload.Location
	event.EventDate = payload.EventDate
	event.Capacity = payload.Capacity

	if err := s.repo.Create(ctx, event); err != nil {
		return err
	}

	return nil
}

func (s *EventService) Update(ctx context.Context, payload models.RequestUpdateEvent) error {
	var event models.Event
	event.Title = payload.Title
	event.Description = payload.Description
	event.Location = payload.Location
	event.EventDate = payload.EventDate
	event.Capacity = payload.Capacity

	if err := s.repo.Update(ctx, event); err != nil {
		return err
	}

	return nil
}

func (s *EventService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *EventService) List(ctx context.Context, filter models.RequestFilterEvent) ([]models.Event, *models.Meta, error) {
	events, meta, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return events, meta, nil
}

func (s *EventService) GetByID(ctx context.Context, id string) (*models.Event, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return event, nil
}

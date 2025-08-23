package handler

import (
	"context"
	"event-booking-system/internal/domain"
	"event-booking-system/internal/models"
	"event-booking-system/internal/service"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{
		service: service,
	}
}

func (h *BookingHandler) Create(c echo.Context) error {
	var payload models.RequestBooking
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	claims := c.Get("user").(*domain.UserClaims)
	payload.UserID = claims.ID

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	if err := h.service.Create(ctx, payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, payload)
}

func (h *BookingHandler) Cancel(c echo.Context) error {
	var payload models.RequestCancelBooking
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	claims := c.Get("user").(*domain.UserClaims)
	payload.UserID = claims.ID

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	if err := h.service.Cancel(ctx, payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "")
}

func (h *BookingHandler) Confirm(c echo.Context) error {
	var payload models.RequestConfirmBooking
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad request",
		})
	}

	claims := c.Get("user").(*domain.UserClaims)
	payload.UserID = claims.ID

	ctx, cancel := context.WithTimeout(c.Request().Context(), 1*time.Second)
	defer cancel()

	if err := h.service.Confim(ctx, payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "")
}

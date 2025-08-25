package main

import (
	"context"
	"event-booking-system/internal/db"
	"event-booking-system/internal/handler"
	"event-booking-system/internal/repository"
	"event-booking-system/internal/router"
	"event-booking-system/internal/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := db.NewDB()
	if err != nil {
		panic(err)
	}

	var greeting string
	if err = db.QueryRow(context.Background(), "SELECT 'Hello, DB!'").Scan(&greeting); err != nil {
		fmt.Printf("QueryRow failed: %v\n", err)
	}
	fmt.Print(greeting)

	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	userService := service.NewUserService(userRepo)
	eventService := service.NewEventService(eventRepo, bookingRepo)
	bookingService := service.NewBookingService(bookingRepo, eventRepo)

	userHandler := handler.NewUserHandler(userService)
	eventHandler := handler.NewEventHandler(eventService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!")
	})

	router.SetupRoutes(e, userHandler, eventHandler, bookingHandler)

	e.Logger.Fatal(e.Start(":8079"))
}

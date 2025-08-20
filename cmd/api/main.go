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
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)
	eventHandler := handler.NewEventHandler(eventService)
	bookingRepo := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(bookingRepo, eventRepo)
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

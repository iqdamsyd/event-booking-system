package main

import (
	"context"
	"event-booking-system/internal/databases/postgres"
	"event-booking-system/internal/databases/redis"
	"event-booking-system/internal/handler"
	middlewares "event-booking-system/internal/middleware"
	"event-booking-system/internal/repository"
	"event-booking-system/internal/router"
	"event-booking-system/internal/service"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := postgres.NewDB()
	if err != nil {
		panic(err)
	}

	var greeting string
	if err = db.QueryRow(context.Background(), "SELECT 'Hello, DB!'").Scan(&greeting); err != nil {
		fmt.Printf("QueryRow failed: %v\n", err)
	}
	fmt.Println(greeting)

	cache, err := redis.NewRedis()
	if err != nil {
		panic(err)
	}

	if err = cache.Set(context.Background(), "foo", "bar", 0).Err(); err != nil {
		fmt.Printf("Cache SET failed: %v\n", err)
	}
	val, err := cache.Get(context.Background(), "foo").Result()
	if err != nil {
		fmt.Printf("Cache GET failed: %v\n", err)
	}
	fmt.Println("foo", val)

	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	userService := service.NewUserService(userRepo, cache)
	eventService := service.NewEventService(eventRepo, bookingRepo)
	bookingService := service.NewBookingService(bookingRepo, eventRepo)

	userHandler := handler.NewUserHandler(userService)
	eventHandler := handler.NewEventHandler(eventService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.RateLimiterMiddleware(cache))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!")
	})

	router.SetupRoutes(e, cache, userHandler, eventHandler, bookingHandler)

	e.Logger.Fatal(e.Start(":8079"))
}

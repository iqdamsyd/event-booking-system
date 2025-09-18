package router

import (
	"event-booking-system/internal/handler"
	"event-booking-system/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

var (
	admin = "admin"
	user  = "user"
)

func SetupRoutes(e *echo.Echo, cache *redis.Client, userHandler *handler.UserHandler, eventHandler *handler.EventHandler, bookingHandler *handler.BookingHandler) {
	private := e.Group("/api", middleware.JWTMiddleware(cache))

	e.POST("/users/register", userHandler.Register)
	e.POST("/users/login", userHandler.Login)

	private.GET("/users", userHandler.List, middleware.RoleMiddleware(admin))
	private.GET("/users/:id", userHandler.GetByID, middleware.RoleMiddleware(admin))
	private.POST("/users", userHandler.Create, middleware.RoleMiddleware(admin))
	private.PUT("/users/:id", userHandler.Update, middleware.RoleMiddleware(admin))
	private.DELETE("/users/:id", userHandler.Delete, middleware.RoleMiddleware(admin))

	private.GET("/events", eventHandler.List)
	private.GET("/events/:id", eventHandler.GetByID)
	private.GET("/events/:id/overview", eventHandler.GetOverview, middleware.RoleMiddleware(admin))
	private.POST("/events", eventHandler.Create, middleware.RoleMiddleware(admin))
	private.PUT("/events/:id", eventHandler.Update, middleware.RoleMiddleware(admin))
	private.DELETE("/events/:id", eventHandler.Delete, middleware.RoleMiddleware(admin))

	private.POST("/bookings", bookingHandler.Create, middleware.RoleMiddleware(user))
	private.PUT("/bookings/cancel", bookingHandler.Cancel, middleware.RoleMiddleware(user))
	private.PUT("/bookings/confirm", bookingHandler.Confirm, middleware.RoleMiddleware(user))
	private.GET("/bookings/my", bookingHandler.List, middleware.RoleMiddleware(user))
}

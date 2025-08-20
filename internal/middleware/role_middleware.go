package middleware

import (
	"event-booking-system/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RoleMiddleware(requiredRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*domain.UserClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing user context")
			}

			for _, role := range requiredRoles {
				if user.Role == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "forbidden")
		}
	}
}

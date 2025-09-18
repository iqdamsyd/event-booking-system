package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/labstack/echo/v4"
)

var (
	rateLimit = 10
)

func RateLimiterMiddleware(rdb *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIP := c.RealIP()
			currentMinute := time.Now().UTC().Minute()
			key := fmt.Sprintf("%s:%d", userIP, currentMinute)

			val, err := rdb.Get(c.Request().Context(), key).Result()
			if err != nil && err != redis.Nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "error getting data from redis")
			}

			if val != "" {
				count, _ := strconv.Atoi(val)
				if count >= rateLimit {
					return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
				}
			}

			_, err = rdb.TxPipelined(c.Request().Context(), func(pipe redis.Pipeliner) error {
				pipe.Incr(c.Request().Context(), key)
				pipe.Expire(c.Request().Context(), key, 59*time.Second)
				return nil
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "error executing redis transaction")
			}

			return next(c)
		}
	}
}

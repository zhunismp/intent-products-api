package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

func AccessLogMiddleware(log *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		log.Info(fmt.Sprintf("%s - %s", c.Method(), c.Path()),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.String("ip", c.IP()),
			slog.Int64("duration", time.Since(start).Milliseconds()),
			slog.Any("request_id", c.Locals("request_id")),
		)

		return err
	}
}

package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	slogctx "github.com/veqryn/slog-context"
)

func RequestIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {

		reqID := uuid.NewString()

		ctx := slogctx.Append(c.Context(),
			"request_id", reqID,
		)

		c.Locals("request_id", reqID)
		c.SetContext(ctx)

		return c.Next()
	}
}

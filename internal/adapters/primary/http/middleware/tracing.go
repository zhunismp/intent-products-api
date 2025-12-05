package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

func TraceMiddleware() fiber.Handler {
	propagator := otel.GetTextMapPropagator()
	tracer := otel.Tracer("http.server")

	return func(c fiber.Ctx) error {
		ctx := propagator.Extract(c.Context(), propagation.HeaderCarrier(c.GetReqHeaders()))
		ctx, span := tracer.Start(ctx, c.Method()+" "+c.Path(),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPRequestMethodOriginal(c.Method()),
				semconv.URLFull(c.OriginalURL()),
				semconv.URLPath(string(c.Request().URI().Path())),
				semconv.ServerAddress(c.Hostname()),
				attribute.String("component", "http-handler"),
			),
		)
		defer span.End()

		start := time.Now()
		c.SetContext(ctx)

		err := c.Next()

		duration := time.Since(start).Seconds()
		statusCode := c.Response().StatusCode()

		span.SetAttributes(
			semconv.HTTPResponseStatusCode(statusCode),
			attribute.Float64("http.server.duration_seconds", duration),
		)

		if err != nil {
			span.RecordError(err)
		}

		return err
	}
}

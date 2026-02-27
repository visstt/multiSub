package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Logger — middleware для структурированного логирования запросов
func Logger(log *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()

		attrs := []any{
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"ip", c.IP(),
			"request_id", c.Locals("requestid"),
		}

		switch {
		case status >= 500:
			log.Error("запрос обработан с ошибкой сервера", attrs...)
		case status >= 400:
			log.Warn("запрос обработан с клиентской ошибкой", attrs...)
		default:
			log.Info("запрос обработан", attrs...)
		}

		return err
	}
}

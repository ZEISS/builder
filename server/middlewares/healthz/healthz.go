package healthz

import "github.com/gofiber/fiber/v3"

// New creates a new middleware handler.
func New() fiber.Handler {
	// Return new handler
	return func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}
}

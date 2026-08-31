package healthz

import "github.com/gofiber/fiber/v3"

// New creates a new middleware handler.
//
// filesystem does not handle url encoded values (for example spaces)
// on it's own. If you need that functionality, set "UnescapePath"
// in fiber.Config.
func New() fiber.Handler {
	// Return new handler
	return func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}
}

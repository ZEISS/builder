package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
)

// SessionController handles session-related requests.
type SessionController interface {
	// GetAccount returns the account associated with the session.
	GetAccount(ctx context.Context) (*models.Account, error)
}

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

// SessionService is the service for managing sessions.
type SessionService interface {
	// GetSession returns the session associated with the request.
	GetSession(ctx context.Context, models *models.Session) error
	// CreateSession creates a new session.
	CreateSession(ctx context.Context) (*models.Session, error)
	// UpdateSession updates an existing session.
	UpdateSession(ctx context.Context, session *models.Session) error
	// DeleteSession deletes a session.
	DeleteSession(ctx context.Context, session *models.Session) error
	// ListSessions returns the list of sessions.
	ListSessions(ctx context.Context, sessions *[]*models.Session) error
}

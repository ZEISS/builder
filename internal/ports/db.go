package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
)

// ReadTx is a read-only transaction interface.
type ReadTx interface {
	// GetCurrentAccount retrieves the current account from the database.
	GetCurrentAccount(ctx context.Context, account *models.Account) error
	// GetSession retrieves the session from the database.
	GetSession(ctx context.Context, session *models.Session) error
	// ListSessions retrieves the list of sessions from the database.
	ListSessions(ctx context.Context, sessions *[]*models.Session) error
}

// ReadWriteTx is a read-write transaction interface.
type ReadWriteTx interface {
	// CreateAccount creates a new account in the database.
	CreateAccount(ctx context.Context, account *models.Account) error
	// CreateSession creates a new session in the database.
	CreateSession(ctx context.Context, session *models.Session) error
	// UpdateSession updates an existing session in the database.
	UpdateSession(ctx context.Context, session *models.Session) error
	// DeleteSession deletes a session in the database.
	DeleteSession(ctx context.Context, session *models.Session) error

	ReadTx
}

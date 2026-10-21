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
	// GetAccount returns the account associated with the session.
	GetAccount(ctx context.Context, account *models.Account) error
	// UpdateAccount updates the account associated with the session.
	UpdateAccount(ctx context.Context, account *models.Account) error
	// DeleteAccount deletes the account associated with the session.
	DeleteAccount(ctx context.Context, account *models.Account) error
	// CreateAccount creates a new account.
	CreateAccount(ctx context.Context, account *models.Account) error
	// ListAccounts returns the list of accounts.
	ListAccounts(ctx context.Context, accounts *[]*models.Account) error
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
	// GetSite returns the site associated with the session.
	GetSite(ctx context.Context, site *models.Site) error
	// CreateSite creates a new site.
	CreateSite(ctx context.Context, site *models.Site) error
	// UpdateSite updates an existing site.
	UpdateSite(ctx context.Context, site *models.Site) error
	// DeleteSite deletes a site.
	DeleteSite(ctx context.Context, site *models.Site) error
}

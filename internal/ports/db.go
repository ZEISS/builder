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
	// GetAccount retrieves the account from the database.
	GetAccount(ctx context.Context, account *models.Account) error
	// ListAccounts retrieves the list of accounts from the database.
	ListAccounts(ctx context.Context, accounts *[]*models.Account) error
	// GetSite retrieves the site from the database.
	GetSite(ctx context.Context, site *models.Site) error
}

// ReadWriteTx is a read-write transaction interface.
type ReadWriteTx interface {
	// CreateSession creates a new session in the database.
	CreateSession(ctx context.Context, session *models.Session) error
	// UpdateSession updates an existing session in the database.
	UpdateSession(ctx context.Context, session *models.Session) error
	// DeleteSession deletes a session in the database.
	DeleteSession(ctx context.Context, session *models.Session) error
	// CreateAccount creates a new account in the database.
	CreateAccount(ctx context.Context, account *models.Account) error
	// UpdateAccount updates an existing account in the database.
	UpdateAccount(ctx context.Context, account *models.Account) error
	// DeleteAccount deletes an account in the database.
	DeleteAccount(ctx context.Context, account *models.Account) error
	// CreateSite creates a new site in the database.
	CreateSite(ctx context.Context, site *models.Site) error
	// UpdateSite updates an existing site in the database.
	UpdateSite(ctx context.Context, site *models.Site) error
	// DeleteSite deletes a site in the database.
	DeleteSite(ctx context.Context, site *models.Site) error

	ReadTx
}

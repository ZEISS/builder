package sessions

import (
	"context"
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/oidc"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/internal/ui/status"
	"gorm.io/gorm"

	"github.com/zeiss/pkg/dbx"
)

// Session is the interface for a session.
type Session interface {
	// ID returns the session ID.
	ID() string
	// Account returns the account associated with the session.
	Account() *models.Account
	// CreateAccount creates a new account for the session.
	CreateAccount(ctx context.Context, token *oidc.Token) (*models.Account, error)
	// HasAccount returns true if the session has an account.
	HasAccount() bool
	// Initialize initializes the session.
	Initialize(ctx context.Context) error
}

// SessionImpl is the implementation of the Session interface.
type SessionImpl struct {
	store   dbx.Database[ports.ReadTx, ports.ReadWriteTx]
	account *models.Account
}

// New returns a new SessionImpl.
func New(store dbx.Database[ports.ReadTx, ports.ReadWriteTx]) Session {
	return &SessionImpl{
		store: store,
	}
}

// ID returns the session ID.
func (s *SessionImpl) ID() string {
	if s.account == nil {
		return ""
	}

	return s.account.ID.String()
}

// CreateAccount creates a new account for the session.
func (s *SessionImpl) CreateAccount(ctx context.Context, token *oidc.Token) (*models.Account, error) {
	account := &models.Account{}
	err := s.store.ReadWriteTx(ctx, func(ctx context.Context, tx ports.ReadWriteTx) error {
		return tx.CreateAccount(ctx, account)
	})
	if err != nil {
		return nil, err
	}

	s.account = account
	return account, nil
}

// GetAccount returns the account associated with the session.
func (s *SessionImpl) Initialize(ctx context.Context) error {
	account := &models.Account{}
	err := s.store.ReadTx(ctx, func(ctx context.Context, tx ports.ReadTx) error {
		return tx.GetCurrentAccount(ctx, account)
	})
	notfound := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !notfound {
		return err
	}

	if notfound {
		return nil
	}

	s.account = account
	return nil
}

// Account returns the account associated with the session.
func (s *SessionImpl) Account() *models.Account {
	return s.account
}

// HasAccount returns true if the session has an account.
func (s *SessionImpl) HasAccount() bool {
	return s.account != nil
}

// CreateAccount creates a new account for the session.
func CreateAccount(session Session, token *oidc.Token) tea.Cmd {
	return func() tea.Msg {
		_, err := session.CreateAccount(context.Background(), token)
		if err != nil {
			return status.Msg(status.Error(err))
		}

		return status.Msg(status.Success("Account created successfully"))
	}
}

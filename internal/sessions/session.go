package sessions

import (
	"context"
	"errors"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"gorm.io/gorm"

	"github.com/zeiss/pkg/dbx"
	"github.com/zeiss/pkg/utilx"
)

// Session is the interface for a session.
type Session interface {
	// Account returns the account associated with the session.
	Account() *models.Account
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

// GetAccount returns the account associated with the session.
func (s *SessionImpl) Initialize(ctx context.Context) error {
	account := &models.Account{}
	err := s.store.ReadTx(ctx, func(ctx context.Context, tx ports.ReadTx) error {
		return tx.GetCurrentAccount(ctx, account)
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
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
	return utilx.NotNil(s.account)
}

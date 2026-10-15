package controllers

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"

	"github.com/zeiss/pkg/dbx"
)

var _ ports.SessionController = (*SessionController)(nil)

// SessionController handles session-related requests.
type SessionController struct {
	store dbx.Database[ports.ReadTx, ports.ReadWriteTx]
}

// NewSessionController returns a new SessionController.
func NewSessionController(store dbx.Database[ports.ReadTx, ports.ReadWriteTx]) *SessionController {
	return &SessionController{
		store: store,
	}
}

// GetAccount returns the account associated with the session.
func (c *SessionController) GetAccount(ctx context.Context) (*models.Account, error) {
	account := &models.Account{}
	err := c.store.ReadTx(ctx, func(ctx context.Context, tx ports.ReadTx) error {
		return tx.GetCurrentAccount(ctx, account)
	})
	if err != nil {
		return nil, err
	}

	return account, nil
}

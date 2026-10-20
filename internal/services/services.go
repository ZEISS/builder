package services

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/pkg/dbx"
)

// Service is the main service struct.
type Service interface {
	// Sessions is the session service.
	ports.SessionService
}

// ServiceImpl is the implementation of the Service interface.
type ServiceImpl struct {
	store dbx.Database[ports.ReadTx, ports.ReadWriteTx]
}

// NewService creates a new service.
func NewService(store dbx.Database[ports.ReadTx, ports.ReadWriteTx]) Service {
	return &ServiceImpl{
		store: store,
	}
}

// CreateSession creates a new session.
func (s *ServiceImpl) CreateSession(ctx context.Context) (*models.Session, error) {
	session := &models.Session{}

	err := s.store.ReadWriteTx(ctx, func(ctx context.Context, rw ports.ReadWriteTx) error {
		return rw.CreateSession(ctx, session)
	})

	return session, err
}

// GetSession returns the session associated with the request.
func (s *ServiceImpl) GetSession(ctx context.Context, session *models.Session) error {
	return s.store.ReadTx(ctx, func(ctx context.Context, tx ports.ReadTx) error {
		return tx.GetSession(ctx, session)
	})
}

// UpdateSession updates the session associated with the request.
func (s *ServiceImpl) UpdateSession(ctx context.Context, session *models.Session) error {
	return s.store.ReadWriteTx(ctx, func(ctx context.Context, rw ports.ReadWriteTx) error {
		return rw.UpdateSession(ctx, session)
	})
}

// DeleteSession deletes the session associated with the request.
func (s *ServiceImpl) DeleteSession(ctx context.Context, session *models.Session) error {
	return s.store.ReadWriteTx(ctx, func(ctx context.Context, rw ports.ReadWriteTx) error {
		return rw.DeleteSession(ctx, session)
	})
}

// ListSessions returns the list of sessions.
func (s *ServiceImpl) ListSessions(ctx context.Context, sessions *[]*models.Session) error {
	return s.store.ReadTx(ctx, func(ctx context.Context, tx ports.ReadTx) error {
		return tx.ListSessions(ctx, sessions)
	})
}

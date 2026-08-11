package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
)

// ReadTx is a read-only transaction interface.
type ReadTx interface {
	Get(ctx context.Context, account *models.Account) error
}

// ReadWriteTx is a read-write transaction interface.
type ReadWriteTx interface{}

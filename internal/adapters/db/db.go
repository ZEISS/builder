package db

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"

	"github.com/google/uuid"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/zeiss/pkg/dbx"
	"github.com/zeiss/pkg/filex"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	_ dbx.Migrator      = (*Database)(nil)
	_ ports.ReadTx      = (*readTxImpl)(nil)
	_ ports.ReadWriteTx = (*writeTxImpl)(nil)
)

// CustomDriverName is the name of the custom SQLite driver.
const CustomDriverName = "sqlite3_extended"

func init() {
	sql.Register(CustomDriverName,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				err := conn.RegisterFunc(
					"gen_random_uuid",
					func(arguments ...any) (string, error) {
						return uuid.New().String(), nil
					},
					true,
				)
				return err
			},
		},
	)
}

type readTxImpl struct {
	conn *gorm.DB
}

// NewReadTx is a factory function that returns a new readTxImpl.
func NewReadTx() dbx.ReadTxFactory[ports.ReadTx] {
	return func(db *gorm.DB) (ports.ReadTx, error) {
		return &readTxImpl{conn: db}, nil
	}
}

// Get implements the ReadTx interface.
func (r *readTxImpl) GetCurrentAccount(ctx context.Context, account *models.Account) error {
	return r.conn.WithContext(ctx).First(account).Error
}

// GetSession implements the ReadTx interface.
func (r *readTxImpl) GetSession(ctx context.Context, session *models.Session) error {
	return r.conn.WithContext(ctx).First(session).Error
}

// ListSessions implements the ReadTx interface.
func (r *readTxImpl) ListSessions(ctx context.Context, sessions *[]*models.Session) error {
	return r.conn.WithContext(ctx).Find(sessions).Error
}

type writeTxImpl struct {
	conn *gorm.DB
	readTxImpl
}

// NewReadWriteTx is a factory function that returns a new writeTxImpl.
func NewWriteTx() dbx.ReadWriteTxFactory[ports.ReadWriteTx] {
	return func(db *gorm.DB) (ports.ReadWriteTx, error) {
		return &writeTxImpl{conn: db}, nil
	}
}

// CreateAccount implements the AccountRepository interface.
func (tx *writeTxImpl) CreateAccount(ctx context.Context, account *models.Account) error {
	return tx.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "current", "refresh_token", "access_token", "expires_at", "token_type", "scope", "id_token", "session_state", "updated_at"}),
		}).Create(account).Error
	})
}

// CreateSession implements the SessionRepository interface.
func (tx *writeTxImpl) CreateSession(ctx context.Context, session *models.Session) error {
	return tx.conn.WithContext(ctx).Create(session).Error
}

// UpdateSession implements the SessionRepository interface.
func (tx *writeTxImpl) UpdateSession(ctx context.Context, session *models.Session) error {
	return tx.conn.WithContext(ctx).Save(session).Error
}

// DeleteSession implements the SessionRepository interface.
func (tx *writeTxImpl) DeleteSession(ctx context.Context, session *models.Session) error {
	return tx.conn.WithContext(ctx).Delete(session).Error
}

// Database is a struct that represents the database adapter.
type Database struct {
	gorm *gorm.DB
	conn *sql.DB
}

var _ ports.AccountRepository = (*Database)(nil)

// Opt is an function that returns the options for the database.
type Opt func(*Database) error

// WithSqlite3 returns an option that sets the database to use SQLite3.
func WithSqlite3(path string) Opt {
	return func(db *Database) error {
		apath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		if err := filex.MkdirAll(apath, 0o700); err != nil {
			return fmt.Errorf("failed to create directory for database: %w", err)
		}

		conn, err := sql.Open(CustomDriverName, filepath.Join(apath, "session.db"))
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		db.conn = conn

		g, err := gorm.Open(sqlite.Dialector{
			DriverName: CustomDriverName,
			DSN:        apath,
			Conn:       conn,
		}, &gorm.Config{
			SkipDefaultTransaction:   true,
			DisableNestedTransaction: true,
		})
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}

		db.gorm = g

		return nil
	}
}

// RunMigrations is a helper function to run the migrations for the database.
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Account{},
		&models.Session{},
	)
}

// New creates a new instance of the Database adapter.
func New(gorm *gorm.DB) *Database {
	return &Database{
		gorm: gorm,
	}
}

// NewDB creates a new instance of the Database adapter.
func NewDB(opts ...Opt) (*Database, error) {
	db := &Database{}

	for _, opt := range opts {
		if err := opt(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// Conn returns the underlying GORM DB connection.
func (d *Database) Conn() *gorm.DB {
	return d.gorm
}

// Migrate is a method that runs the migrations for the database.
func (d *Database) Migrate(ctx context.Context, _ ...any) error {
	return d.gorm.WithContext(ctx).AutoMigrate(&models.Account{}, &models.Session{})
}

// Get is a method that returns an account by ID.
func (d *Database) Get(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).First(account).Error
}

// GetCurrent is a method that returns the current account.
func (d *Database) GetCurrent(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).First(account).Where(&models.Account{Current: true}).Error
}

// Create is a method that creates a new account.
func (d *Database) Create(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "current", "refresh_token", "access_token", "expires_at", "token_type", "scope", "id_token", "session_state", "updated_at"}),
	}).Create(account).Error
}

// Update is a method that updates an account.
func (d *Database) Update(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).Save(account).Error
}

// Delete is a method that deletes an account.
func (d *Database) Delete(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).Delete(account).Error
}

// List is a method that returns a list of all accounts.
func (d *Database) List(ctx context.Context, accounts *[]models.Account) error {
	return d.gorm.WithContext(ctx).Find(accounts).Error
}

// Close is a method that closes the database connection.
func (d *Database) Close() error {
	return d.conn.Close()
}

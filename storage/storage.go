package storage

import (
	"database/sql"
	"fmt"

	"github.com/crypto-trade/config"
	_ "github.com/lib/pq"
)

const (
	dbManufacturer = "postgres"
)

// Storage represents storage structure.
type Storage struct {
	db *sql.DB
}

// New returns new storage implementation.
func New(config config.DB) (*Storage, error) {
	db, err := sql.Open(dbManufacturer, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name),
	)
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		db: db,
	}

	return storage, nil
}

// Close closes DB.
func (s *Storage) Close() error {
	return s.db.Close()
}

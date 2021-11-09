package storage

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

const (
	duplicateKeyValueViolatesUniqueConstraint pq.ErrorCode = "23505"
)

// Storage errors.
var (
	ErrNotFound                                  error = errors.New("not found")
	ErrDuplicateKeyValueViolatesUniqueConstraint error = errors.New("pq: duplicate key value violates unique constraint")
)

func errorFromCode(pqError error) error {
	if pqError == nil {
		return nil
	}

	if errors.Is(pqError, sql.ErrNoRows) {
		return ErrNotFound
	}

	e, ok := pqError.(*pq.Error)
	if !ok {
		return pqError
	}

	switch e.Code {
	case duplicateKeyValueViolatesUniqueConstraint:
		return ErrDuplicateKeyValueViolatesUniqueConstraint
	}

	return e
}

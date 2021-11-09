package repository

import (
	"context"

	"github.com/crypto-trade/model"
)

// User describes user interface.
type User interface {
	AddNewUser(ctx context.Context, user model.User) (model.User, error)
	HashedPassword(ctx context.Context, email model.Email) (model.Password, error)
}

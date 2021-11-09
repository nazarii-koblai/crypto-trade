package model

import "github.com/gofrs/uuid"

// Wallet represents wallet.
type Wallet struct {
	ID       uuid.UUID `json:"id"`
	Currency string    `json:"currency"`
	Balance  int       `json:"balance"`
}

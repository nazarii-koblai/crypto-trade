package model

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	minLength = 8
	maxLength = 50
)

// Password represents password.
type Password string

// Validate validates password.
func (p Password) Validate() error {
	if len(p) < minLength || len(p) > maxLength {
		return errors.New("password length should be in range [8-50] symbols")
	}
	return nil
}

// Hash hashes password using bcrypt.
func (p *Password) Hash() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(string(*p)), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("can't hash password: %w", err)
	}
	*p = Password(hash)
	return nil
}

// CompareWithHashed compares password with hashed password.
func (p *Password) CompareWithHashed(hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(*p))
}

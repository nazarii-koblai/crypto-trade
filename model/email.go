package model

import (
	"fmt"
	"net"
	"net/mail"
	"strings"
)

// Email represents email struct.
type Email string

// Validates an email.
func (e Email) Validate() error {
	if _, err := mail.ParseAddress(string(e)); err != nil {
		return fmt.Errorf("%s isn't a valid email address", e)
	}
	domain := strings.Split(string(e), "@")[1]
	if _, err := net.LookupMX(domain); err != nil {
		return fmt.Errorf("%s has invalid domain: %w", e, err)
	}
	return nil
}

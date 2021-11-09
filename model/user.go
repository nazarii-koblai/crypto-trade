package model

import "fmt"

// User respresents user struct.
type User struct {
	Name     PersonalInfo `json:"username"`
	Surname  PersonalInfo `json:"surname"`
	Email    Email        `json:"email"`
	Password Password     `json:"password"`
	Wallets  []Wallet     `json:"wallets"`
}

// Validate validates user.
func (u User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("username is empty")
	}

	if u.Surname == "" {
		return fmt.Errorf("surname is empty")
	}

	for _, validator := range []interface{ Validate() error }{
		u.Name,
		u.Surname,
		u.Email,
		u.Password,
	} {
		if err := validator.Validate(); err != nil {
			return err
		}
	}
	return nil
}

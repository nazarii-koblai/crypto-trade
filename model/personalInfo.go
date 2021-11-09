package model

import (
	"fmt"
	"unicode"
)

// PersonalInfo represents personal info.
type PersonalInfo string

// Validate validates personal info.
func (pi PersonalInfo) Validate() error {
	for _, s := range pi {
		if !unicode.IsLetter(s) {
			return fmt.Errorf("%s should contain only unicode letters", pi)
		}
	}
	return nil
}

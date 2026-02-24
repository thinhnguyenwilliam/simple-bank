package util

import "net/mail"

// IsValidEmail checks whether the email format is valid
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

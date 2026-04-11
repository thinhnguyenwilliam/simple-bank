// simple-bank/val/validator.go
package val

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[\p{L} .'-]+$`).MatchString // hỗ trợ Unicode (tên VN)
)

// ValidateEmail validates email format
func ValidateEmail(value string) error {
	if err := ValidateString(value, "email", 320); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("email is not a valid email address")
	}
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(value string) error {
	if err := ValidateString(value, "password", 128); err != nil {
		return err
	}
	if len(value) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, r := range value {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("password must include upper, lower, number and special character")
	}
	return nil
}

// ValidateFullName validates person's full name
func ValidateFullName(value string) error {
	if err := ValidateString(value, "full_name", 100); err != nil {
		return err
	}
	if len(value) < 2 {
		return fmt.Errorf("full_name must be at least 2 characters")
	}
	if !isValidFullName(value) {
		return fmt.Errorf("full_name contains invalid characters")
	}
	return nil
}

// ValidateString validates a generic string
func ValidateString(value string, name string, maxLength int) error {
	if value == "" {
		return fmt.Errorf("%s must not be empty", name)
	}

	if len(value) > maxLength {
		return fmt.Errorf("%s must not be longer than %d characters", name, maxLength)
	}

	return nil
}

// ValidateUsername validates username rules
func ValidateUsername(value string) error {
	if err := ValidateString(value, "username", 50); err != nil {
		return err
	}

	if len(value) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}

	if !isValidUsername(value) {
		return fmt.Errorf("username must contain only letters, numbers or underscores")
	}

	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	return ValidateString(value, "secret_code", 128)
}

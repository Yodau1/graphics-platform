package users_service

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	core_errors "github.com/Yodau1/graphics-platform/internal/core/errors"
)

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 || len(password) > 72 {
		return fmt.Errorf(
			"password must contain at least 8 characters and at most 72 bytes: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	if err := validatePassword(password); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return string(hash), nil
}

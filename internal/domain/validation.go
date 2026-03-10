package domain

import (
	"fmt"
	"strings"
)

func ValidateCreateInput(key, value string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("%w: key is required", ErrValidation)
	}
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w: value is required", ErrValidation)
	}
	return nil
}

func ValidateDeleteInput(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}
	return nil
}

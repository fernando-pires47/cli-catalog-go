package domain

import "errors"

var (
	ErrNotFound       = errors.New("command not found")
	ErrValidation     = errors.New("validation error")
	ErrAmbiguous      = errors.New("ambiguous command match")
	ErrDangerDenied   = errors.New("dangerous command denied")
	ErrInvalidCatalog = errors.New("invalid catalog json")
)

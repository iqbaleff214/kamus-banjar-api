package contribution

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("contribution not found")
	ErrWordNotFound  = errors.New("word not found")
	ErrForbidden     = errors.New("not allowed to modify this submission")
	ErrWordActive    = errors.New("approved words cannot be modified or deleted")
	ErrNotPending    = errors.New("only pending submissions can be deleted")
	ErrDuplicateWord = errors.New("word already exists")
)

// ErrInvalidReq wraps a validation message as an error.
func ErrInvalidReq(msg string) error {
	return fmt.Errorf("%w: %s", errValidation, msg)
}

// errValidation is the sentinel used by ErrInvalidReq.
var errValidation = errors.New("invalid request")

// IsValidationErr reports whether err is a validation error from ErrInvalidReq.
func IsValidationErr(err error) bool {
	return errors.Is(err, errValidation)
}

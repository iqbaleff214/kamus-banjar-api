package user

import "errors"

var (
	ErrEmailTaken      = errors.New("email already registered")
	ErrInvalidCreds    = errors.New("invalid email or password")
	ErrAccountInactive = errors.New("account is deactivated")
	ErrTokenInvalid    = errors.New("refresh token is invalid or expired")
	ErrWrongPassword   = errors.New("current password is incorrect")
	ErrNotFound        = errors.New("user not found")
	ErrSelfModify      = errors.New("cannot modify your own account")
)

package postgres

import "errors"

// Errors
var (
	ErrInvalidDBSecret = errors.New("Failed to get a valid vault postgres secret")
	ErrClosed          = errors.New("Client has been closed")
	ErrInvalidAddr     = errors.New("Address cannot be empty")
	ErrInvalidDB       = errors.New("DB cannot be empty")
)

package database

import (
	"errors"
)

var (
	// ErrNotFound AS not found.
	ErrNotFound = errors.New("AS not found")
)

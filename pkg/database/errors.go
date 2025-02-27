package database

import (
	"errors"
)

var (
	// ErrASNotFound is returned when an autonomous system for the specified IP was not found.
	ErrASNotFound = errors.New("autonomous system not found")
)

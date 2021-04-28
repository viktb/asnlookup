package binarytrie

import "errors"

var (
	// ErrInvalidIPAddress IP address is invalid.
	ErrInvalidIPAddress = errors.New("invalid IP address")
	// ErrTrieImmutable Trie is immutable.
	ErrTrieImmutable = errors.New("trie is immutable")
	// ErrValueNotFound value was not found.
	ErrValueNotFound = errors.New("value not found")
	// ErrInvalidFormat invalid marshaled input.
	ErrInvalidFormat = errors.New("invalid format")
)

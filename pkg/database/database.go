package database

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/viktb/asnlookup/pkg/binarytrie"
)

// AutonomousSystem represents an Autonomous System on the Internet.
type AutonomousSystem struct {
	// Number (aka ASN) is the unique identifier for an Autonomous System.
	Number uint32
}

// Database stores mappings between IP addresses and Autonomous Systems.
type Database struct {
	mappings *binarytrie.ArrayTrie
}

// NewFromDump creates a Database from a binary dump.
func NewFromDump(r io.Reader) (*Database, error) {
	d := &Database{
		mappings: binarytrie.NewArrayTrie(),
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %v", err)
	}

	if err = d.UnmarshalBinary(data); err != nil {
		return nil, fmt.Errorf("failed to restore dump: %v", err)
	}
	return d, nil
}

// Lookup returns the AutonomousSystem for a given net.IP.
func (d *Database) Lookup(ip net.IP) (AutonomousSystem, error) {
	asn, err := d.mappings.Lookup(ip)
	if errors.Is(err, binarytrie.ErrValueNotFound) {
		return AutonomousSystem{}, ErrASNotFound
	} else if err != nil {
		return AutonomousSystem{}, fmt.Errorf("lookup failed: %w", err)
	}

	return AutonomousSystem{
		Number: asn,
	}, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (d *Database) MarshalBinary() ([]byte, error) {
	return d.mappings.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (d *Database) UnmarshalBinary(data []byte) error {
	return d.mappings.UnmarshalBinary(data)
}

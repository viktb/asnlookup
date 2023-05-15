package database

import (
	"encoding"
	"fmt"
	"io"
	"io/ioutil"
	"net"

	"github.com/banviktor/asnlookup/pkg/binarytrie"
)

// AutonomousSystem represents an Autonomous System on the Internet.
type AutonomousSystem struct {
	// Number (aka ASN) is the unique identifier for an Autonomous System.
	Number uint32
}

// Database stores mappings between IP addresses and Autonomous Systems.
type Database interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	// Lookup returns the AutonomousSystem for a given net.IP.
	Lookup(net.IP) (AutonomousSystem, error)
}

type database struct {
	mappings *binarytrie.ArrayTrie
}

// Lookup implements Database.
func (d *database) Lookup(ip net.IP) (AutonomousSystem, error) {
	asn, err := d.mappings.Lookup(ip)
	if err == binarytrie.ErrValueNotFound {
		return AutonomousSystem{}, ErrNotFound
	} else if err != nil {
		return AutonomousSystem{}, fmt.Errorf("lookup failed: %v", err)
	}

	return AutonomousSystem{
		Number: asn,
	}, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (d *database) MarshalBinary() ([]byte, error) {
	return d.mappings.MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (d *database) UnmarshalBinary(data []byte) error {
	return d.mappings.UnmarshalBinary(data)
}

func NewFromDump(r io.Reader) (Database, error) {
	d := &database{
		mappings: binarytrie.NewArrayTrie(),
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %v", err)
	}

	if err = d.UnmarshalBinary(data); err != nil {
		return nil, fmt.Errorf("failed to restore dump: %v", err)
	}
	return d, nil
}

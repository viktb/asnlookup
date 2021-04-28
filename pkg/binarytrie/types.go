package binarytrie

import (
	"net"
)

type Trie interface {
	// Insert inserts an IP network - value mapping into the trie.
	Insert(*net.IPNet, uint32) error
	// Lookup returns a value for the given IP address.
	Lookup(ip net.IP) (uint32, error)
}

package database

import (
	"fmt"
	"io"
	"net"

	"github.com/banviktor/go-mrt"

	"github.com/banviktor/asnlookup/pkg/binarytrie"
)

type builder struct {
	prototype  *binarytrie.NaiveTrie
	fillFactor float32
}

// InsertMapping stores an IP prefix - AutonomousSystem mapping.
func (b *builder) InsertMapping(ipNet *net.IPNet, asn uint32) error {
	err := b.prototype.Insert(ipNet, asn)
	if err != nil {
		return err
	}
	return nil
}

// ImportMRT imports records from an MRT stream.
func (b *builder) ImportMRT(input io.Reader) error {
	r := mrt.NewReader(input)

	for {
		record, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to parse MRT record: %v", err)
		}

		rib, ok := record.(*mrt.TableDumpV2RIB)
		if !ok || isNullMask(rib.Prefix.Mask) {
			continue
		}

		prefix, asn, err := mrtRIBToMapping(rib)
		if err != nil {
			continue
		}

		err = b.InsertMapping(prefix, asn)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetFillFactor sets the fill factor parameter for the optimization phase.
func (b *builder) SetFillFactor(fillFactor float32) {
	b.fillFactor = fillFactor
}

// Build builds the Database instance.
func (b *builder) Build() (Database, error) {
	err := b.prototype.Optimize(b.fillFactor)
	if err != nil {
		return nil, err
	}
	return &database{
		mappings: b.prototype.ToArrayTrie(),
	}, nil
}

// NewBuilder creates a builder.
func NewBuilder() *builder {
	return &builder{
		prototype:  binarytrie.NewNaiveTrie(),
		fillFactor: 0.5,
	}
}

func isNullMask(mask net.IPMask) bool {
	for _, b := range mask {
		if b != 0 {
			return false
		}
	}
	return true
}

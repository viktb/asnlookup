package database

import (
	"fmt"
	"io"
	"net"

	"github.com/viktb/go-mrt"

	"github.com/viktb/asnlookup/pkg/binarytrie"
)

// Builder is a helper for constructing a Database.
type Builder struct {
	prototype  *binarytrie.NaiveTrie
	fillFactor float32
}

// NewBuilder creates a Builder.
func NewBuilder() *Builder {
	return &Builder{
		prototype:  binarytrie.NewNaiveTrie(),
		fillFactor: 0.5,
	}
}

// InsertMapping stores an IP prefix - AutonomousSystem mapping.
func (b *Builder) InsertMapping(ipNet *net.IPNet, asn uint32) error {
	err := b.prototype.Insert(ipNet, asn)
	if err != nil {
		return err
	}
	return nil
}

// ImportMRT imports records from an MRT stream.
func (b *Builder) ImportMRT(input io.Reader) error {
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
func (b *Builder) SetFillFactor(fillFactor float32) {
	b.fillFactor = fillFactor
}

// Build builds the Database instance.
func (b *Builder) Build() (*Database, error) {
	err := b.prototype.Optimize(b.fillFactor)
	if err != nil {
		return nil, err
	}
	return &Database{
		mappings: b.prototype.ToArrayTrie(),
	}, nil
}

func isNullMask(mask net.IPMask) bool {
	for _, b := range mask {
		if b != 0 {
			return false
		}
	}
	return true
}

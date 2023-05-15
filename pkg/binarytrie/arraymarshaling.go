package binarytrie

import (
	"encoding/binary"
)

const (
	arrayTrieMarshalHeader  = "github.com/banviktor/asnlookup/pkg/binarytrie\x00ArrayTrie\x00"
	arrayTrieMarshalVersion = uint8(1)
)

// MarshalBinary implements encoding.BinaryMarshaler.
func (t *ArrayTrie) MarshalBinary() (data []byte, err error) {
	var i int
	data = make([]byte, len(arrayTrieMarshalHeader)+1+8+len(t.nodes)*10+8+len(t.skippedBits)*12)

	// Write header.
	copy(data, arrayTrieMarshalHeader)
	i += len(arrayTrieMarshalHeader)
	data[i] = arrayTrieMarshalVersion
	i += 1

	// Write nodes.
	binary.LittleEndian.PutUint64(data[i:i+8], uint64(len(t.nodes)))
	i += 8
	for _, n := range t.nodes {
		bN, err := n.MarshalBinary()
		if err != nil {
			return nil, err
		}
		copy(data[i:i+10], bN)
		i += 10
	}

	// Write skipped bits information.
	binary.LittleEndian.PutUint64(data[i:i+8], uint64(len(t.skippedBits)))
	i += 8
	for k, v := range t.skippedBits {
		binary.LittleEndian.PutUint64(data[i:i+8], uint64(k))
		i += 8
		binary.LittleEndian.PutUint32(data[i:i+4], v)
		i += 4
	}

	return
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *ArrayTrie) UnmarshalBinary(data []byte) error {
	var i int

	// Check header.
	if string(data[:len(arrayTrieMarshalHeader)]) != arrayTrieMarshalHeader {
		return ErrInvalidFormat
	}
	i += len(arrayTrieMarshalHeader)
	if data[i] != arrayTrieMarshalVersion {
		return ErrInvalidFormat
	}
	i += 1

	// Populate nodes.
	nodeCount := binary.LittleEndian.Uint64(data[i : i+8])
	i += 8
	t.nodes = make([]arrayTrieNode, nodeCount)
	for j := uint64(0); j < nodeCount; j++ {
		if err := t.nodes[j].UnmarshalBinary(data[i : i+10]); err != nil {
			return err
		}
		i += 10
	}

	// Populate skipped bits information.
	skippedBitCount := binary.LittleEndian.Uint64(data[i : i+8])
	i += 8
	t.skippedBits = make(map[int]uint32, skippedBitCount)
	for j := uint64(0); j < skippedBitCount; j++ {
		k := int(binary.LittleEndian.Uint64(data[i : i+8]))
		i += 8
		v := binary.LittleEndian.Uint32(data[i : i+4])
		i += 4
		t.skippedBits[k] = v
	}

	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (n *arrayTrieNode) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 10)
	data[0] = n.branchingFactor
	data[1] = n.skipValue
	binary.LittleEndian.PutUint32(data[2:6], n.childrenOffset)
	binary.LittleEndian.PutUint32(data[6:10], n.value)
	return
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (n *arrayTrieNode) UnmarshalBinary(data []byte) error {
	n.branchingFactor = data[0]
	n.skipValue = data[1]
	n.childrenOffset = binary.LittleEndian.Uint32(data[2:6])
	n.value = binary.LittleEndian.Uint32(data[6:10])
	return nil
}

package binarytrie

import (
	"fmt"
	"net"
)

// ArrayTrie represents a trie in a space-efficient array format.
type ArrayTrie struct {
	nodes       []arrayTrieNode
	skippedBits map[int]uint32
}

type arrayTrieNode struct {
	branchingFactor uint8
	skipValue       uint8
	childrenOffset  uint32
	value           uint32
}

// Insert implements Trie.
func (t *ArrayTrie) Insert(*net.IPNet, uint32) error {
	return ErrTrieImmutable
}

// Lookup implements Trie.
func (t *ArrayTrie) Lookup(ip net.IP) (value uint32, err error) {
	var bitPosition, index, nextIndex int
	var skippedBits uint32

	for {
		if t.nodes[index].value != 0 {
			value = t.nodes[index].value
		}
		if t.nodes[index].isLeaf() {
			break
		}

		skippedBits = extractBits(ip, bitPosition, int(t.nodes[index].skipValue))
		bitPosition += int(t.nodes[index].skipValue)
		nextIndex = index + int(t.nodes[index].childrenOffset) + int(extractBits(ip, bitPosition, int(t.nodes[index].branchingFactor)))
		bitPosition += int(t.nodes[index].branchingFactor)

		if t.nodes[index].skipValue > 0 {
			if expected, ok := t.skippedBits[nextIndex]; !ok || expected != skippedBits {
				break
			}
		}
		index = nextIndex
	}
	if value == 0 {
		return 0, ErrValueNotFound
	}
	return
}

// NewArrayTrie returns an empty ArrayTrie.
//
// This is rarely useful, as an ArrayTrie does not support inserts.
func NewArrayTrie() *ArrayTrie {
	return &ArrayTrie{
		nodes:       make([]arrayTrieNode, 0),
		skippedBits: make(map[int]uint32),
	}
}

// NewArrayTrieFromNaiveTrie creates an ArrayTrie from a NaiveTrie.
func NewArrayTrieFromNaiveTrie(nt *NaiveTrie) *ArrayTrie {
	at := NewArrayTrie()
	at.nodes = make([]arrayTrieNode, 0, nt.allocatedSize())

	nodeQueue := []*naiveTrieNode{nt.root}
	for len(nodeQueue) > 0 {
		batchIndex := len(at.nodes)
		batchSize := len(nodeQueue)
		for _, nNode := range nodeQueue {
			if nNode == nil {
				at.nodes = append(at.nodes, arrayTrieNode{})
				continue
			}

			i := len(at.nodes)
			at.nodes = append(at.nodes, arrayTrieNode{
				branchingFactor: nNode.branchingFactor,
				skipValue:       nNode.skipValue,
				value:           nNode.value,
			})
			if nNode.parent != nil && nNode.parent.skipValue > 0 {
				at.skippedBits[i] = nNode.skippedBits
			}
			if !nNode.isLeaf() {
				at.nodes[i].childrenOffset = uint32(batchIndex - i + len(nodeQueue))
				nodeQueue = append(nodeQueue, nNode.children...)
			}
		}
		nodeQueue = nodeQueue[batchSize:]
	}

	return at
}

func (n *arrayTrieNode) isLeaf() bool {
	return n.branchingFactor == 0
}

// String implements fmt.Stringer.
func (t *ArrayTrie) String() string {
	str := "#\tBF\tSV\tCO\tValue\n"
	for i, n := range t.nodes {
		suffix := ""
		if bits, ok := t.skippedBits[i]; ok {
			suffix = fmt.Sprintf(" (skipped: %0*b)", n.skipValue, bits)
		}
		str += fmt.Sprintf("%d\t%d\t%d\t%d\t%d%s\n", i, n.branchingFactor, n.skipValue, n.childrenOffset, n.value, suffix)
	}
	str += fmt.Sprintf("%v\n", t.skippedBits)
	return str
}

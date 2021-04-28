package binarytrie

import (
	"net"
)

type NaiveTrie struct {
	root    *naiveTrieNode
	mutable bool
}

type naiveTrieNode struct {
	skipValue       uint8
	skippedBits     uint32
	branchingFactor uint8
	parent          *naiveTrieNode
	children        []*naiveTrieNode
	value           uint32
}

// NewNaiveTrie creates an empty NaiveTrie.
func NewNaiveTrie() *NaiveTrie {
	return &NaiveTrie{
		root:    &naiveTrieNode{},
		mutable: true,
	}
}

// Insert implements Trie.
func (t *NaiveTrie) Insert(ipNet *net.IPNet, value uint32) error {
	if !t.mutable {
		return ErrTrieImmutable
	}

	prefix, prefixSize, err := parseIpNet(ipNet)
	if err != nil {
		return err
	}

	currentNode := t.root
	bitPosition := 0
	for {
		if currentNode.branchingFactor == 0 {
			currentNode.branchingFactor = 1
			currentNode.children = make([]*naiveTrieNode, 2)
		}

		bit := extractBits(prefix, bitPosition, 1)
		bitPosition++
		if currentNode.children[bit] == nil {
			currentNode.children[bit] = &naiveTrieNode{parent: currentNode}
		}
		currentNode = currentNode.children[bit]

		if bitPosition >= prefixSize {
			break
		}
	}
	currentNode.value = value
	return nil
}

// Lookup implements Trie.
func (t *NaiveTrie) Lookup(ip net.IP) (value uint32, err error) {
	ip = ip.To16()
	if ip == nil {
		return 0, ErrInvalidIPAddress
	}

	currentNode := t.root
	bitPosition := 0
	for {
		if currentNode.value != 0 {
			value = currentNode.value
		}
		if currentNode.isLeaf() {
			break
		}

		skippedBits := extractBits(ip, bitPosition, int(currentNode.skipValue))
		bitPosition += int(currentNode.skipValue)
		prefix := extractBits(ip, bitPosition, int(currentNode.branchingFactor))
		bitPosition += int(currentNode.branchingFactor)

		nextNode := currentNode.children[prefix]
		if nextNode == nil || nextNode.skippedBits != skippedBits {
			break
		}
		currentNode = nextNode
	}
	if value == 0 {
		return 0, ErrValueNotFound
	}
	return
}

// ToArrayTrie creates an identical ArrayTrie.
func (t *NaiveTrie) ToArrayTrie() *ArrayTrie {
	return NewArrayTrieFromNaiveTrie(t)
}

func (t *NaiveTrie) allocatedSize() int {
	return t.root.allocatedSize()
}

func (n *naiveTrieNode) isLeaf() bool {
	return n.branchingFactor == 0
}

func (n *naiveTrieNode) allocatedSize() int {
	count := 1
	for _, child := range n.children {
		if child != nil {
			count += child.allocatedSize()
		} else {
			count++
		}
	}
	return count
}

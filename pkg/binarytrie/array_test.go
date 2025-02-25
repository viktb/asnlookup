package binarytrie_test

import (
	"testing"
)

func TestEmptyArrayTrieLookup(t *testing.T) {
	trie, testCases := newEmptyNaiveTrie()
	testLookup(t, trie.ToArrayTrie(), testCases)
}

func TestTrivialArrayTrieLookup(t *testing.T) {
	trie, testCases := newTrivialNaiveTrie(t)
	testLookup(t, trie.ToArrayTrie(), testCases)
}

func TestPopulatedArrayTrieLookup(t *testing.T) {
	trie, testCases := newPopulatedNaiveTrie(t)
	testLookup(t, trie.ToArrayTrie(), testCases)
}

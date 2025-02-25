package binarytrie_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/viktb/asnlookup/pkg/binarytrie"
)

func TestEmptyMarshaledArrayTrieLookup(t *testing.T) {
	trie, testCases := newEmptyNaiveTrie()
	assert.NoError(t, trie.Optimize(0.5), "Optimize should not error")
	arrayTrie := trie.ToArrayTrie()

	buf, err := arrayTrie.MarshalBinary()
	assert.NoError(t, err, "MarshalBinary should not error")

	newTrie := &ArrayTrie{}
	err = newTrie.UnmarshalBinary(buf)
	assert.NoError(t, err, "UnmarshalBinary should not error")

	testLookup(t, newTrie, testCases)
}

func TestTrivialMarshaledArrayTrieLookup(t *testing.T) {
	trie, testCases := newTrivialNaiveTrie(t)
	assert.NoError(t, trie.Optimize(0.5), "Optimize should not error")
	arrayTrie := trie.ToArrayTrie()

	buf, err := arrayTrie.MarshalBinary()
	assert.NoError(t, err, "MarshalBinary should not error")

	newTrie := &ArrayTrie{}
	err = newTrie.UnmarshalBinary(buf)
	assert.NoError(t, err, "UnmarshalBinary should not error")

	testLookup(t, newTrie, testCases)
}

func TestPopulatedMarshaledArrayTrieLookup(t *testing.T) {
	trie, testCases := newPopulatedNaiveTrie(t)
	assert.NoError(t, trie.Optimize(0.5), "Optimize should not error")
	arrayTrie := trie.ToArrayTrie()

	buf, err := arrayTrie.MarshalBinary()
	assert.NoError(t, err, "MarshalBinary should not error")

	newTrie := &ArrayTrie{}
	err = newTrie.UnmarshalBinary(buf)
	assert.NoError(t, err, "UnmarshalBinary should not error")

	testLookup(t, newTrie, testCases)
}

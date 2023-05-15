package binarytrie_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyOptimizedNaiveTrieLookup(t *testing.T) {
	trie, testCases := newEmptyNaiveTrie()
	assert.NoError(t, trie.Optimize(0.5), "Optimize should not error")
	testLookup(t, trie, testCases)
}

func TestTrivialOptimizedNaiveTrieLookup(t *testing.T) {
	trie, testCases := newTrivialNaiveTrie()
	assert.NoError(t, trie.Optimize(0.5), "Optimize should not error")
	testLookup(t, trie, testCases)
}

func TestPopulatedOptimizedNaiveTrieLookup(t *testing.T) {
	trie, testCases := newPopulatedNaiveTrie()
	assert.NoError(t, trie.Optimize(0.5), "Optimize should not error")
	testLookup(t, trie, testCases)
}

func TestEmptyOptimizedArrayTrieLookup(t *testing.T) {
	trie, testCases := newEmptyNaiveTrie()
	assert.NoError(t, trie.Optimize(0.5), "Optimize should not error")
	arrayTrie := trie.ToArrayTrie()

	fmt.Println(arrayTrie.String())
	testLookup(t, arrayTrie, testCases)
}

func TestTrivialOptimizedArrayTrieLookup(t *testing.T) {
	trie, testCases := newTrivialNaiveTrie()
	assert.NoError(t, trie.Optimize(0.5), "Optimize should not error")
	arrayTrie := trie.ToArrayTrie()

	fmt.Println(arrayTrie.String())
	testLookup(t, arrayTrie, testCases)
}

func TestPopulatedOptimizedArrayTrieLookup(t *testing.T) {
	trie, testCases := newPopulatedNaiveTrie()
	assert.NoError(t, trie.Optimize(0.5), "Optimize should not error")
	arrayTrie := trie.ToArrayTrie()

	fmt.Println(arrayTrie.String())
	testLookup(t, arrayTrie, testCases)
}

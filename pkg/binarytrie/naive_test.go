package binarytrie_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/banviktor/asnlookup/pkg/binarytrie"
)

type testCase struct {
	ip  string
	asn uint32
	err error
}

func TestEmptyNaiveTrieLookup(t *testing.T) {
	trie, testCases := newEmptyNaiveTrie()
	testLookup(t, trie, testCases)
}

func TestTrivialNaiveTrieLookup(t *testing.T) {
	trie, testCases := newTrivialNaiveTrie()
	testLookup(t, trie, testCases)
}

func TestPopulatedNaiveTrieLookup(t *testing.T) {
	trie, testCases := newPopulatedNaiveTrie()
	testLookup(t, trie, testCases)
}

func newEmptyNaiveTrie() (*NaiveTrie, []testCase) {
	trie := NewNaiveTrie()

	testCases := []testCase{
		{"0.0.0.0", 0, ErrValueNotFound},
		{"255.255.255.255", 0, ErrValueNotFound},
	}

	return trie, testCases
}

func newTrivialNaiveTrie() (*NaiveTrie, []testCase) {
	trie := NewNaiveTrie()
	_, ipNet, _ := net.ParseCIDR("0.0.0.0/0")
	trie.Insert(ipNet, 42)

	testCases := []testCase{
		{"0.0.0.0", 42, nil},
		{"255.255.255.255", 42, nil},
	}

	return trie, testCases
}

func newPopulatedNaiveTrie() (*NaiveTrie, []testCase) {
	trie := NewNaiveTrie()
	testData := []struct {
		net string
		asn uint32
	}{
		{"192.168.1.0/24", 999},
		{"0.0.0.0/2", 200},
		{"128.0.0.0/2", 210},
		{"160.0.0.0/3", 2101},
		{"160.0.0.0/3", 2101}, // duplicate entry on purpose
		{"192.0.0.0/3", 211},
		{"224.0.0.0/3", 211},
	}
	for _, td := range testData {
		_, ipNet, _ := net.ParseCIDR(td.net)
		trie.Insert(ipNet, td.asn)
	}

	testCases := []testCase{
		{"0.0.0.0", 200, nil},
		{"32.128.128.128", 200, nil},
		{"63.255.255.255", 200, nil},
		{"64.0.0.0", 0, ErrValueNotFound},
		{"96.128.128.128", 0, ErrValueNotFound},
		{"127.255.255.255", 0, ErrValueNotFound},
		{"128.0.0.0", 210, nil},
		{"159.255.255.255", 210, nil},
		{"160.128.128.128", 2101, nil},
		{"191.255.255.255", 2101, nil},
		{"192.0.0.0", 211, nil},
		{"192.168.0.255", 211, nil},
		{"192.168.1.0", 999, nil},
		{"192.168.1.128", 999, nil},
		{"192.168.1.255", 999, nil},
		{"192.168.2.0", 211, nil},
		{"224.128.128.128", 211, nil},
		{"255.255.255.255", 211, nil},
	}

	return trie, testCases
}

func testLookup(t *testing.T, trie Trie, testCases []testCase) {
	for _, tc := range testCases {
		asn, err := trie.Lookup(net.ParseIP(tc.ip))
		if tc.err != nil && assert.Error(t, err, "%s should have error", tc.ip) {
			assert.Equal(t, tc.err, err)
		}
		assert.Equal(t, int(tc.asn), int(asn), "%s expected AS%d, actual: AS%d", tc.ip, tc.asn, asn)
	}
}

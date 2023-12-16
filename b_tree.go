package main

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// immutable b+tree

type BNode struct {
	data []byte // byte array, that will be dumped to the disk
}

const (
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 2 // leaf nodes with values
)

type BTree struct {
	//pointer (nonzero page number)
	root uint64

	// callbacks for managing on-disk pages
	get func(uint64) BNode // dereference a pointer
	new func(BNode) uint64 // allocate a new page
	del func(uint64)       // deallocate a page
}

const HEADER = 4
const BTREE_PAGE_SIZE = 4096 // 4K bytes
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

var Assert *assert.Assertions

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	t := &testing.T{}
	Assert = assert.New(t)
	assert.True(t, node1max <= BTREE_PAGE_SIZE, "node1max cannot be greater than BTREE_PAGE_SIZE!")
}

// header
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], btype)
}

// pointers
func (node BNode) getPtr(idx uint16, t *testing.T) uint64 {
	// assert(idx < node.nkeys())
	condition := idx < node.nkeys()
	assert.True(t, condition, "index is not less than node.nkeys()!")
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64, t *testing.T) {
	condition := idx < node.nkeys()
	assert.True(t, condition, "index is not less than node.nkeys()!")
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

// offset list

func offsetPos(node BNode, idx uint16, t *testing.T) uint16 {
	condition := 1 <= idx && idx <= node.nkeys()
	assert.True(t, condition, "index out of bounds!")
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx, nil):]) // not sure if `nil` will work (?)
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx, nil):], offset) // (?)
}

// key-values

func (node BNode) kvPos(idx uint16, t *testing.T) uint16 {
	condition := idx <= node.nkeys()
	assert.True(t, condition, "index cannot be greater than number of nodes!")
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16, t *testing.T) []byte {
	condition := idx < node.nkeys()
	assert.True(t, condition, "index out of bounds!")

	pos := node.kvPos(idx, nil) // (?)
	klen := binary.LittleEndian.Uint16(node.data[pos:])
	return node.data[pos+4:][:klen]
}

func (node BNode) getVal(idx uint16, t *testing.T) []byte {
	condition := idx < node.nkeys()
	assert.True(t, condition, "index out of bounds!")

	pos := node.kvPos(idx, nil) // (?)
	klen := binary.LittleEndian.Uint16(node.data[pos+0:])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])
	return node.data[pos+4+klen:][:vlen]
}

// node size in bytes
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys(), nil) // (?)
}

/*
B-Tree Insertion
*/

// returns the first child node whose range intersects the key. (child[i] <= key)
// TODO: bisect
func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)

	// the first key is a copy from the parent node
	// thus its always less than or equal to the key

	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i, nil), key) // (?)
		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}
	return found
}

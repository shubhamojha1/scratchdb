package main

import (
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

func offsetPos(node BNode, idx uint16, t *testing.T) {
	condition := 1 <= idx && idx <= node.nkeys()
	assert.True(t, condition, "index out of bounds!")
}

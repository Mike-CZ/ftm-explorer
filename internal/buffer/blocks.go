package buffer

import (
	"fmt"
	"ftm-explorer/internal/types"
)

// BlocksBuffer represents a buffer of blocks. It is used to store blocks
// in the order of their numbers. The buffer is cyclic, i.e. when the buffer
// is full, the oldest block is removed from the buffer and the new block
// is added to the buffer. The buffer relies on the fact that the block
// numbers are monotonically increasing. That means that the block with
// the number N+1 is always added to the buffer after the block with the
// number N. The buffer is not thread-safe.
type BlocksBuffer struct {
	// data is a slice of blocks
	data []*types.Block
	// capacity is the capacity of the buffer
	capacity uint
	// head is the index of the last inserted block
	head uint
	// size is the number of blocks in the buffer
	size uint
}

// NewBlocksBuffer creates a new blocks buffer.
// The size of the buffer is specified by the size parameter.
func NewBlocksBuffer(size uint) *BlocksBuffer {
	return &BlocksBuffer{
		data:     make([]*types.Block, size),
		capacity: size,
		head:     0,
		size:     0,
	}
}

// Get returns the block with the specified number.
// If the block is not found, the second return value is false.
func (bb *BlocksBuffer) Get(number uint64) (*types.Block, bool) {
	blk := bb.data[bb.index(number)]

	// if the block is not found, return nil
	if blk == nil {
		return nil, false
	}

	// if the block is found on given index, check the number
	if uint64(blk.Number) != number {
		return nil, false
	}

	return blk, true
}

// GetLatest returns the latest inserted blocks.
// The number of blocks to return is specified by the number parameter.
// If the number of blocks in the buffer is less than the number parameter,
// all blocks in the buffer are returned.
func (bb *BlocksBuffer) GetLatest(count uint) []*types.Block {
	// if the number of blocks in the buffer is less than the number of
	// blocks to return, return all blocks
	if bb.size < count {
		count = bb.size
	}

	// create a slice of blocks
	blocks := make([]*types.Block, count)

	// copy the blocks to the slice
	for i := uint(0); i < count; i++ {
		blocks[i] = bb.data[(bb.head-i+bb.capacity)%bb.capacity]
	}

	return blocks
}

// Add adds the specified block to the buffer.
func (bb *BlocksBuffer) Add(block *types.Block) {
	index := bb.index(uint64(block.Number))

	// increase the size if the block is new
	if bb.data[index] == nil {
		bb.size++
	}

	// add the block to the buffer
	bb.data[index] = block

	// set head to the index of last inserted block
	bb.head = index
}

// Len returns the number of blocks in the buffer.
func (bb *BlocksBuffer) Len() uint {
	return bb.size
}

// index returns the index of the block with the specified number.
func (bb *BlocksBuffer) index(number uint64) uint {
	return uint(number % uint64(bb.capacity))
}

// printBuffer prints the buffer to the standard output.
func (bb *BlocksBuffer) printBuffer() {
	for i := 0; i < len(bb.data); i++ {
		blk := bb.data[i]
		if blk == nil {
			continue
		}
		fmt.Printf("%d: %d\n", i, uint64(blk.Number))
	}
}

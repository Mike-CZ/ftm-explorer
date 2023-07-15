package buffer

import (
	"ftm-explorer/internal/types"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Test empty buffer has length 0 and no blocks
func TestBlocksBuffer_Empty(t *testing.T) {
	bb := NewBlocksBuffer(3)

	// check the buffer is empty
	if bb.Len() != 0 {
		t.Error("expected length to be 0")
	}

	_, ok := bb.Get(1)
	if ok {
		t.Error("expected block to not be found")
	}
}

// Test block to buffer can be added and retrieved
func TestBlocksBuffer_AddAndGet(t *testing.T) {
	bb := NewBlocksBuffer(5)

	blocksNumbers := []uint64{3, 4, 5, 6}

	for _, number := range blocksNumbers {
		// add block
		bb.Add(&types.Block{Number: hexutil.Uint64(number)})
		// check if block is returned
		value, ok := bb.Get(number)
		if !ok {
			t.Fatalf("expected block number %d to be found", number)
		}
		if uint64(value.Number) != number {
			t.Errorf("expected block number %d to be returned, got %d", number, uint64(value.Number))
		}
	}

	// assert length is 4
	if bb.Len() != 4 {
		t.Error("expected length to be 4")
	}
}

// Test latest blocks can be retrieved
func TestBlocksBuffer_GetLatest(t *testing.T) {
	bb := NewBlocksBuffer(10)

	// define 10 block numbers
	blocksNumbers := []uint64{14, 15, 16, 17, 18, 19, 20, 21, 22, 23}

	for _, number := range blocksNumbers {
		bb.Add(&types.Block{Number: hexutil.Uint64(number)})
	}

	// check if latest 5 blocks are returned
	retrieved := bb.GetLatest(5)
	if len(retrieved) != 5 {
		t.Error("expected 5 blocks to be returned")
	}

	// check if blocks are returned in correct order
	for i, block := range retrieved {
		if uint64(block.Number) != blocksNumbers[len(blocksNumbers)-i-1] {
			t.Error("expected blocks to be returned in correct order")
		}
	}

	// getting more blocks than inserted should return all blocks
	retrieved = bb.GetLatest(20)
	if len(retrieved) != 10 {
		t.Error("expected 10 blocks to be returned")
	}

	// assert length is 10
	if bb.Len() != 10 {
		t.Error("expected length to be 10")
	}
}

// Test blocks are overwritten when buffer is full
func TestBlocksBuffer_BlocksAreRewritten(t *testing.T) {
	bb := NewBlocksBuffer(5)

	// add 10 blocks
	blocksNumbers := []uint64{14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	for _, number := range blocksNumbers {
		bb.Add(&types.Block{Number: hexutil.Uint64(number)})
	}

	// assert first 5 blocks are not present in buffer
	for i := uint64(0); i < 5; i++ {
		_, ok := bb.Get(blocksNumbers[i])
		if ok {
			t.Error("expected block to not be found")
		}
	}

	// assert last 5 blocks are present in buffer
	for i := uint64(5); i < 10; i++ {
		value, ok := bb.Get(blocksNumbers[i])
		if !ok {
			t.Fatalf("expected block to be found")
		}
		if uint64(value.Number) != blocksNumbers[i] {
			t.Errorf("expected block '%d' to be returned, got %d", blocksNumbers[i], uint64(value.Number))
		}
	}

	// assert length is 5
	if bb.Len() != 5 {
		t.Error("expected length to be 5")
	}
}

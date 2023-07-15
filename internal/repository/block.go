package repository

import (
	"context"
	"ftm-explorer/internal/types"

	eth "github.com/ethereum/go-ethereum/core/types"
)

// GetBlockByNumber returns the block identified by number.
// If the block is not in the buffer, it will be fetched from the RPC.
func (r *Repository) GetBlockByNumber(number uint64) (*types.Block, error) {
	// try to get block from buffer
	blk, exists := r.blkBuffer.Get(number)
	if exists {
		return blk, nil
	}

	// get block from rpc if not exists in buffer
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()
	blk, err := r.rpc.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	return blk, nil
}

// GetLatestObservedBlocks returns the number of latest observed blocks.
// It will only return blocks that are in the buffer.
func (r *Repository) GetLatestObservedBlocks(count uint) []*types.Block {
	return r.blkBuffer.GetLatest(count)
}

// GetLatestObservedBlock returns the latest observed block.
// It will only return a block that is in the buffer (if any).
func (r *Repository) GetLatestObservedBlock() *types.Block {
	blocks := r.blkBuffer.GetLatest(1)
	if len(blocks) == 0 {
		return nil
	}
	return blocks[0]
}

// UpdateLatestObservedBlock updates the latest observed block.
// It will add the block to the buffer.
func (r *Repository) UpdateLatestObservedBlock(blk *types.Block) {
	r.blkBuffer.Add(blk)
}

// GetNewHeadersChannel returns a channel that will receive the latest headers from blockchain.
func (r *Repository) GetNewHeadersChannel() <-chan *eth.Header {
	return r.rpc.ObservedHeadProxy()
}

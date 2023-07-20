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
func (r *Repository) UpdateLatestObservedBlock(blk *types.Block) error {
	// add block to buffer
	r.blkBuffer.Add(blk)

	// add block to db
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()
	return r.db.AddBlock(ctx, blk)
}

// GetNewHeadersChannel returns a channel that will receive the latest headers from blockchain.
func (r *Repository) GetNewHeadersChannel() <-chan *eth.Header {
	return r.rpc.ObservedHeadProxy()
}

// GetTrxCountAggByTimestamp returns aggregation of transactions in given time range.
// It will fetch data from last block timestamp if endTime is nil.
func (r *Repository) GetTrxCountAggByTimestamp(resolution types.AggResolution, ticks uint, endTime *uint64) ([]types.HexUintTick, error) {
	last := r.getLastBlockTimestamp(endTime)
	if last == nil {
		return nil, nil
	}
	// get aggregation from db
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()
	return r.db.TrxCountAggByTimestamp(ctx, *last, resolution.ToDuration(), ticks)
}

// GetGasUsedAggByTimestamp returns aggregation of gas used in given time range.
// It will fetch data from last block timestamp if endTime is nil.
func (r *Repository) GetGasUsedAggByTimestamp(resolution types.AggResolution, ticks uint, endTime *uint64) ([]types.HexUintTick, error) {
	last := r.getLastBlockTimestamp(endTime)
	if last == nil {
		return nil, nil
	}
	// get aggregation from db
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()
	return r.db.GasUsedAggByTimestamp(ctx, *last, resolution.ToDuration(), ticks)
}

// getLastBlockTimestamp returns the timestamp of the last block if `endTime` is nil.
func (r *Repository) getLastBlockTimestamp(endTime *uint64) *uint64 {
	// if end time is given, then use it
	if endTime != nil {
		return endTime
	}

	// get last observed block
	lastBlock := r.GetLatestObservedBlock()

	// if last observed block is nil, then we can't get endTime
	if lastBlock == nil {
		return nil
	}

	last := uint64(lastBlock.Timestamp)
	return &last
}

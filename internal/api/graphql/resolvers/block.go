package resolvers

import (
	"fmt"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Block represents resolvable blockchain block structure.
type Block struct {
	types.Block
}

// NewBlock builds new resolvable block structure.
func NewBlock(blk *types.Block) *Block {
	if blk == nil {
		return nil
	}
	return &Block{Block: *blk}
}

// RecentBlocks resolves recent observed blocks.
func (rs *RootResolver) RecentBlocks(args *struct{ Limit int32 }) ([]*Block, error) {
	if args.Limit <= 0 {
		return nil, fmt.Errorf("invalid limit value")
	}
	blocks := rs.repository.GetLatestObservedBlocks(int(args.Limit))

	if len(blocks) == 0 {
		return []*Block{}, nil
	}

	rv := make([]*Block, len(blocks))
	for i, b := range blocks {
		rv[i] = NewBlock(b)
	}

	return rv, nil
}

// CurrentBlockHeight resolves current block height.
func (rs *RootResolver) CurrentBlockHeight() (*int32, error) {
	lastBlock := rs.repository.GetLatestObservedBlock()
	if lastBlock == nil {
		return nil, nil
	}
	height := int32(lastBlock.Number)
	return &height, nil
}

// Block resolves block by number.
func (rs *RootResolver) Block(args *struct{ Number hexutil.Uint64 }) (*Block, error) {
	block, err := rs.repository.GetBlockByNumber(uint64(args.Number))
	if err != nil {
		rs.log.Warningf("Failed to get block by number [%d]; %v", args.Number, err)
		return nil, err
	}
	return NewBlock(block), nil
}

// TransactionsCount resolves number of transactions in the block.
func (blk *Block) TransactionsCount() int32 {
	return int32(len(blk.Transactions))
}

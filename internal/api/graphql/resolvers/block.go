package resolvers

import (
	"fmt"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Block represents resolvable blockchain block structure.
type Block struct {
	rs *RootResolver
	types.Block
}

// Tick represents resolvable blockchain tick structure.
type Tick types.HexUintTick

// BlockTimestampAggregations resolves block timestamp aggregations.
func (rs *RootResolver) BlockTimestampAggregations(args *struct{ Subject types.AggSubject }) ([]Tick, error) {
	// get data based on subject
	var result []types.HexUintTick

	switch args.Subject {
	case types.AggSubjectTxsCount:
		result = rs.repository.GetTxCountPer10Secs()
	case types.AggSubjectGasUsed:
		result = rs.repository.GetGasUsedPer10Secs()
	default:
		return nil, fmt.Errorf("invalid subject value")
	}

	// convert result
	rv := make([]Tick, len(result))
	for i, t := range result {
		rv[i] = (Tick)(t)
	}

	return rv, nil
}

// RecentBlocks resolves recent observed blocks.
func (rs *RootResolver) RecentBlocks(args *struct{ Limit int32 }) ([]*Block, error) {
	if args.Limit <= 0 {
		return nil, fmt.Errorf("invalid limit value")
	}
	blocks := rs.repository.GetLatestObservedBlocks(uint(args.Limit))

	if len(blocks) == 0 {
		return []*Block{}, nil
	}

	rv := make([]*Block, len(blocks))
	for i, b := range blocks {
		blk := Block{rs: rs, Block: *b}
		rv[i] = &blk
	}

	return rv, nil
}

// CurrentBlockHeight resolves current block height.
func (rs *RootResolver) CurrentBlockHeight() (*hexutil.Uint64, error) {
	lastBlock := rs.repository.GetLatestObservedBlock()
	if lastBlock == nil {
		return nil, nil
	}
	return &lastBlock.Number, nil
}

// Block resolves block by number.
func (rs *RootResolver) Block(args *struct{ Number hexutil.Uint64 }) (*Block, error) {
	block, err := rs.repository.GetBlockByNumber(uint64(args.Number))
	if err != nil {
		rs.log.Warningf("Failed to get block by number [%d]; %v", args.Number, err)
		return nil, err
	}
	blk := Block{rs: rs, Block: *block}
	return &blk, nil
}

// TransactionsCount resolves number of transactions in the block.
func (blk *Block) TransactionsCount() int32 {
	return int32(len(blk.Transactions))
}

// FullTransactions resolves full transactions in the block.
func (blk *Block) FullTransactions() ([]*Transaction, error) {
	result := make([]*Transaction, 0)

	// fetch transactions
	for _, hash := range blk.Transactions {
		trx, err := blk.rs.repository.GetTransactionByHash(hash)
		if err != nil {
			blk.rs.log.Warningf("Failed to get transaction by hash [%s]; %v", hash.Hex(), err)
			return nil, err
		}
		result = append(result, (*Transaction)(trx))
	}

	return result, nil
}

// Timestamp resolves tick timestamp.
func (t Tick) Timestamp() int32 {
	return int32(t.Time)
}

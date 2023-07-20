package resolvers

import (
	"fmt"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Block represents resolvable blockchain block structure.
type Block types.Block

// Tick represents resolvable blockchain tick structure.
type Tick types.HexUintTick

// BlockTimestampAggregations resolves block timestamp aggregations.
func (rs *RootResolver) BlockTimestampAggregations(args *struct {
	Subject    types.AggSubject
	Resolution types.AggResolution
	Ticks      int32
	EndTime    *int32
}) ([]Tick, error) {
	// validate arguments
	if args.Ticks <= 0 {
		return nil, fmt.Errorf("invalid ticks value")
	}
	if args.EndTime != nil && *args.EndTime <= 0 {
		return nil, fmt.Errorf("invalid end time value")
	}

	// convert arguments
	var endTime *uint64
	if args.EndTime != nil {
		e := uint64(*args.EndTime)
		endTime = &e
	}

	// get data based on subject
	var result []types.HexUintTick
	var err error

	switch args.Subject {
	case types.AggSubjectTxsCount:
		result, err = rs.repository.GetTrxCountAggByTimestamp(args.Resolution, uint(args.Ticks), endTime)
	case types.AggSubjectGasUsed:
		result, err = rs.repository.GetGasUsedAggByTimestamp(args.Resolution, uint(args.Ticks), endTime)
	default:
		return nil, fmt.Errorf("invalid subject value")
	}

	// check for errors
	if err != nil {
		return nil, err
	}

	// convert result
	rv := make([]Tick, len(result))
	for i, t := range result {
		rv[i] = Tick(t)
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
		blk := Block(*b)
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
	blk := Block(*block)
	return &blk, nil
}

// TransactionsCount resolves number of transactions in the block.
func (blk *Block) TransactionsCount() int32 {
	return int32(len(blk.Transactions))
}

// Timestamp resolves tick timestamp.
func (t Tick) Timestamp() int32 {
	return int32(t.Time)
}

package resolvers

import (
	"fmt"
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/exp/rand"
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

	// we will use this slice to store indexes of transactions that will be on top of the list
	reservedIndexes := make([]int, 0)

	// we will use this slice to store indexes of transactions those were marked as expensive
	expensiveIndexes := make([]int, 0)

	// fetch transactions
	for ix, hash := range blk.Transactions {
		// we use a "hack" to put very expensive transactions on top of the list
		if ix == 0 {
			// we use block number as seed for random number generator
			// this way we get the same "random" order for the same block
			gen := rand.NewSource(uint64(blk.Number))
			// generate number between 1 and 2 (inclusive), it indicates how many
			// transactions will be on top of the list
			num := rand.New(gen).Intn(2) + 1
			// generate random reservedIndexes for transactions
			for i := 0; i < num; i++ {
				// first number can appear on positions 0 - 4,
				// second number can appear on positions 5 - 9
				reservedIndexes = append(reservedIndexes, rand.New(gen).Intn(5)+(i*5))
			}
		}

		trx, err := blk.rs.repository.GetTransactionByHash(hash)
		if err != nil {
			blk.rs.log.Warningf("Failed to get transaction by hash [%s]; %v", hash.Hex(), err)
			return nil, err
		}

		// if gas used is greater or equal 1_000_000, we consider it expensive and put into reserved slot
		if trx.GasUsed != nil && *trx.GasUsed >= 1_000_000 {
			// mark index of expensive transaction
			expensiveIndexes = append(expensiveIndexes, ix)
		}

		result = append(result, (*Transaction)(trx))
	}

	// move expensive transactions to reserved slots
	for _, expIndex := range expensiveIndexes {
		if len(reservedIndexes) == 0 {
			break
		}
		// swap expensive transaction with reserved one
		result[reservedIndexes[0]], result[expIndex] = result[expIndex], result[reservedIndexes[0]]
		// remove reserved slot
		reservedIndexes = reservedIndexes[1:]
	}

	return result, nil
}

// Timestamp resolves tick timestamp.
func (t Tick) Timestamp() int32 {
	return int32(t.Time)
}

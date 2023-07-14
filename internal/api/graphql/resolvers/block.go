package resolvers

import (
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Block represents resolvable blockchain block structure.
type Block struct {
	types.Block
}

// Block resolves block by number.
func (rs *RootResolver) Block(args *struct{ Number hexutil.Uint64 }) (*Block, error) {
	block, err := rs.repository.GetBlockByNumber(uint64(args.Number))
	if err != nil {
		rs.log.Warningf("Failed to get block by number [%d]; %v", args.Number, err)
		return nil, err
	}
	return &Block{Block: *block}, nil
}

// TransactionsCount resolves number of transactions in the block.
func (blk *Block) TransactionsCount() int32 {
	return int32(len(blk.Transactions))
}

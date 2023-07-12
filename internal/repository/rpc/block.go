package rpc

import (
	"context"
	"fmt"
	"ftm-explorer/internal/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BlockByNumber returns the block by the given number.
func (rpc *OperaRpc) BlockByNumber(ctx context.Context, number uint64) (*types.Block, error) {
	var block types.Block

	// get the block by number
	err := rpc.ftm.CallContext(ctx, &block, "ftm_getBlockByNumber", hexutil.EncodeUint64(number), false)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %v", err)
	}

	// detect block not found situation; block number is zero and the hash is also zero
	if uint64(block.Number) == 0 && block.Hash.Big().Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("block %d not found", number)
	}

	return &block, nil
}

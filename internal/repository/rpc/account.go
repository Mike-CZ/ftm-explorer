package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// AccountBalance returns the balance of the account.
func (rpc *OperaRpc) AccountBalance(ctx context.Context, addr common.Address) (*hexutil.Big, error) {
	// use RPC to make the call
	var balance string
	err := rpc.ftm.CallContext(ctx, &balance, "eth_getBalance", addr.Hex(), "latest")
	if err != nil {
		return nil, err
	}

	// decode the response from remote server
	val, err := hexutil.DecodeBig(balance)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(val), nil
}

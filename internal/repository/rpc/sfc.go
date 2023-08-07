package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// NumberOfValidators returns the number of validators.
func (rpc *OperaRpc) NumberOfValidators(ctx context.Context) (uint64, error) {
	// get latest sealed epoch
	var epoch hexutil.Bytes
	err := rpc.ftm.CallContext(ctx, &epoch, "ftm_call", map[string]interface{}{
		"to":   rpc.sfcAddress,
		"data": hexutil.Bytes([]byte{0x7c, 0xac, 0xb1, 0xd6})}, "latest")
	if err != nil {
		return 0, err
	}

	// get number of validators in the epoch
	var out hexutil.Bytes
	err = rpc.ftm.CallContext(ctx, &out, "ftm_call", map[string]interface{}{
		"to":   rpc.sfcAddress,
		"data": hexutil.Bytes(append([]byte{0xb8, 0x8a, 0x37, 0xe2}, []byte(epoch)...))}, "latest")
	if err != nil {
		return 0, err
	}

	// the result is offset: 32 bytes, length: 32 bytes
	// so we need to extract the last 32 bytes
	val := new(big.Int).SetBytes(out[32:64])

	return val.Uint64(), nil
}

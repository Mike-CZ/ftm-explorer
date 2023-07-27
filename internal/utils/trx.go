package utils

import (
	"bytes"
	"ftm-explorer/internal/types"
)

const (
	kTrxDeploymentType        = "Deployment"
	kTrxSimpleTxType          = "Simple Tx"
	kTrxErc20TransferType     = "ERC20 Transfer"
	kTrxErc20TransferFromType = "ERC20 TransferFrom"
	kTrxErc20MintType         = "ERC20 Mint"
	kTrxErc20ApproveType      = "ERC20 Approve"
	kTrxOtherTxType           = "Other Tx"
)

// ParseTrxType parses the transaction type from the transaction.
func ParseTrxType(trx *types.Transaction) string {
	// if receiver is empty or null address, it is a deployment
	if trx.To == nil || trx.To.Hex() == "0x0000000000000000000000000000000000000000" {
		return kTrxDeploymentType
	}

	// if data is empty, it is a simple tx
	if len(trx.Input) == 0 {
		return kTrxSimpleTxType
	}

	// if first 4 bytes of data is "0xa9059cbb", it is an erc20 transfer
	if len(trx.Input) >= 4 && bytes.Equal(trx.Input[:4], []byte{0xa9, 0x05, 0x9c, 0xbb}) {
		return kTrxErc20TransferType
	}

	// if first 4 bytes of data is "0x23b872dd", it is an erc20 transferFrom
	if len(trx.Input) >= 4 && bytes.Equal(trx.Input[:4], []byte{0x23, 0xb8, 0x72, 0xdd}) {
		return kTrxErc20TransferFromType
	}

	// if first 4 bytes of data is "0x40c10f19", it is an erc20 mint
	if len(trx.Input) >= 4 && bytes.Equal(trx.Input[:4], []byte{0x40, 0xc1, 0x0f, 0x19}) {
		return kTrxErc20MintType
	}

	// if first 4 bytes of data is "0x095ea7b3", it is an erc20 approve
	if len(trx.Input) >= 4 && bytes.Equal(trx.Input[:4], []byte{0x09, 0x5e, 0xa7, 0xb3}) {
		return kTrxErc20ApproveType
	}

	return kTrxOtherTxType
}

package utils

import (
	"ftm-explorer/internal/types"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestTrx_ParseTrxType(t *testing.T) {
	trx := types.Transaction{}

	// test it is deployment, since the receiver is empty
	if ParseTrxType(&trx) != kTrxDeploymentType {
		t.Errorf("expected type '%s', got %s", kTrxDeploymentType, ParseTrxType(&trx))
	}

	// set the receiver
	receiver := common.HexToAddress("0x1234567890123456789012345678901234567890")
	trx.To = &receiver

	// test it is simple transfer, since the data is empty
	if ParseTrxType(&trx) != kTrxSimpleTxType {
		t.Errorf("expected type '%s', got %s", kTrxSimpleTxType, ParseTrxType(&trx))
	}

	// test erc20 transfer
	trx.Input = []byte{0xa9, 0x05, 0x9c, 0xbb}
	if ParseTrxType(&trx) != kTrxErc20TransferType {
		t.Errorf("expected type '%s', got %s", kTrxErc20TransferType, ParseTrxType(&trx))
	}

	// test erc20 transferFrom
	trx.Input = []byte{0x23, 0xb8, 0x72, 0xdd}
	if ParseTrxType(&trx) != kTrxErc20TransferFromType {
		t.Errorf("expected type '%s', got %s", kTrxErc20TransferFromType, ParseTrxType(&trx))
	}

	// test erc20 mint
	trx.Input = []byte{0x40, 0xc1, 0x0f, 0x19}
	if ParseTrxType(&trx) != kTrxErc20MintType {
		t.Errorf("expected type '%s', got %s", kTrxErc20MintType, ParseTrxType(&trx))
	}

	// test erc20 approve
	trx.Input = []byte{0x09, 0x5e, 0xa7, 0xb3}
	if ParseTrxType(&trx) != kTrxErc20ApproveType {
		t.Errorf("expected type '%s', got %s", kTrxErc20ApproveType, ParseTrxType(&trx))
	}
}

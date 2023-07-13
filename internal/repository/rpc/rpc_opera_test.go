// instruct go test to build this file only when the tag opera_rpc is specified
//go:build opera_rpc
// +build opera_rpc

package rpc

import (
	"context"
	"ftm-explorer/internal/config"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Test that the transaction by hash is returned correctly.
func TestOperaRpc_TransactionByHash(t *testing.T) {
	rpc := createOperaRpc(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// pick random transaction from the blockchain
	// e.g. https://ftmscan.com/tx/0x640beacf2b750682b04faa1e3f8096a524d953dcf11a9f9d4eb55d59323729d7
	trx, err := rpc.TransactionByHash(ctx, common.HexToHash("0x640beacf2b750682b04faa1e3f8096a524d953dcf11a9f9d4eb55d59323729d7"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// check that values are correct
	if trx.Hash != common.HexToHash("0x640beacf2b750682b04faa1e3f8096a524d953dcf11a9f9d4eb55d59323729d7") {
		t.Errorf("unexpected hash: %v", trx.Hash)
	}
	if trx.BlockHash != common.HexToHash("0x0003298d000011e5c9dc093c38330d28c8fa7c42c1494634857f3539f4527e0e") {
		t.Errorf("unexpected block hash: %v", trx.BlockHash)
	}
	if trx.BlockNumber != 61_118_012 {
		t.Errorf("unexpected block number: %v", trx.BlockNumber)
	}
	if trx.From != common.HexToAddress("0x0f83b3f4e728ea15c862498063f29fdb463d043c") {
		t.Errorf("unexpected from address: %v", trx.From)
	}
	if trx.To == nil || *trx.To != common.HexToAddress("0x8f8ddaca443ceac1ee5676721d14cfc5c4548020") {
		t.Errorf("unexpected to address: %v", trx.To)
	}
	if uint64(trx.GasUsed) != 64_848 {
		t.Errorf("unexpected gas used: %v", trx.GasUsed)
	}
	if trx.GasPrice.ToInt().Cmp(big.NewInt(217_949_672_299)) != 0 {
		t.Errorf("unexpected gas price: %v", trx.GasPrice)
	}
	if len(trx.Logs) != 1 {
		t.Errorf("unexpected number of logs: %v", len(trx.Logs))
	}
}

// Test that the block by number is returned correctly.
func TestOperaRpc_BlockByNumber(t *testing.T) {
	rpc := createOperaRpc(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// pick random transaction from the blockchain
	// e.g. https://ftmscan.com/block/61118012
	blk, err := rpc.BlockByNumber(ctx, 61_118_012)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// check that values are correct
	if blk.Number != 61_118_012 {
		t.Errorf("unexpected number: %v", blk.Number)
	}
	if blk.Epoch != 207_245 {
		t.Errorf("unexpected epoch: %v", blk.Epoch)
	}
	if blk.Hash != common.HexToHash("0x0003298d000011e5c9dc093c38330d28c8fa7c42c1494634857f3539f4527e0e") {
		t.Errorf("unexpected hash: %v", blk.Hash)
	}
	if blk.GasUsed != 2_566_683 {
		t.Errorf("unexpected gas used: %v", blk.GasUsed)
	}
	if blk.TimeStamp != 1_682_833_408 {
		t.Errorf("unexpected timestamp: %v", blk.TimeStamp)
	}
	if len(blk.Txs) != 16 {
		t.Errorf("unexpected number of transactions: %v", len(blk.Txs))
	}
}

// Test observed head proxy channel.
func TestOperaRpc_ObservedHeadProxy(t *testing.T) {
	rpc := createOperaRpc(t)

	ch := rpc.ObservedHeadProxy()

	// listen for 5 seconds
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	wg := sync.WaitGroup{}
	wg.Add(1)

	// and count incoming messages
	count := 0
	go func() {
		for {
			select {
			case <-ch:
				count++
			case <-timer.C:
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()

	// validate that we received some messages
	if count == 0 {
		t.Errorf("no messages received")
	}
}

func createOperaRpc(t *testing.T) *OperaRpc {
	rpc, err := NewOperaRpc(&config.Rpc{OperaRpcUrl: "https://rpcapi.fantom.network"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Cleanup(func() {
		rpc.Close()
	})
	return rpc
}

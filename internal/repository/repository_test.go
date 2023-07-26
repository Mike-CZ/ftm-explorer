package repository

import (
	"ftm-explorer/internal/repository/db"
	"ftm-explorer/internal/repository/rpc"
	"ftm-explorer/internal/types"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
)

// Test that repository contains no blocks after initialization.
func TestRepository_NoObservedBlocks(t *testing.T) {
	repository, _, _ := createRepository(t)

	// test latest block is nil
	latestBlock := repository.GetLatestObservedBlock()
	if latestBlock != nil {
		t.Errorf("expected nil, got %v", latestBlock)
	}

	// test latest blocks is empty
	latestBlocks := repository.GetLatestObservedBlocks(10)
	if len(latestBlocks) != 0 {
		t.Errorf("expected empty, got %v", latestBlocks)
	}
}

// Test that repository returns block by number.
func TestRepository_GetAndUpdateBlockByNumber(t *testing.T) {
	repository, mockRpc, mockDb := createRepository(t)

	// block should be returned from rpc because it is not in the buffer
	block := types.Block{Number: 100}

	// expect rpc.BlockByNumber to be called with block number 100
	mockRpc.EXPECT().BlockByNumber(gomock.Any(), gomock.Eq(uint64(100))).Return(&block, nil)
	returnedBlock, err := repository.GetBlockByNumber(100)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if returnedBlock.Number != block.Number {
		t.Errorf("expected %v, got %v", block.Number, returnedBlock.Number)
	}

	// block should be added into database on updating latest observed block
	mockDb.EXPECT().AddBlock(gomock.Any(), gomock.Eq(&block)).Return(nil)

	// update latest observed block
	if err := repository.UpdateLatestObservedBlock(&block); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// block should be returned from buffer because it is in the buffer
	returnedBlock, err = repository.GetBlockByNumber(100)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if returnedBlock.Number != block.Number {
		t.Errorf("expected %v, got %v", block.Number, returnedBlock.Number)
	}
}

// Test that repository returns transaction by hash.
func TestRepository_GetTransactionByHash(t *testing.T) {
	repository, mockRpc, _ := createRepository(t)

	// transaction should be returned from rpc
	trx := types.Transaction{Hash: common.HexToHash("0x123")}
	mockRpc.EXPECT().TransactionByHash(gomock.Any(), gomock.Eq(trx.Hash)).Return(&trx, nil)
	returnedTrx, err := repository.GetTransactionByHash(trx.Hash)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if returnedTrx.Hash != trx.Hash {
		t.Errorf("expected %v, got %v", trx.Hash, returnedTrx.Hash)
	}
}

// Test that repository returns transaction by hash.
func TestRepository_GetNewHeadersChannel(t *testing.T) {
	repository, mockRpc, _ := createRepository(t)

	// transaction should be returned from rpc
	ch := make(chan *eth.Header, 10)

	go func() {
		// send 10 headers into the channel
		for i := 1; i <= 10; i++ {
			ch <- &eth.Header{Number: big.NewInt(int64(i))}
		}
		close(ch)
	}()

	// expect rpc.ObservedHeadProxy to be called
	mockRpc.EXPECT().ObservedHeadProxy().Return(ch)
	rch := repository.GetNewHeadersChannel()

	// read 10 headers from the channel
	number := big.NewInt(1)
	for header := range rch {
		if header.Number.Cmp(number) != 0 {
			t.Errorf("expected %v, got %v", number.Uint64(), header.Number.Uint64())
		}
		number.Add(number, big.NewInt(1))
	}
}

// Test that number of accounts is set and returned correctly.
func TestRepository_GetAndSetNumberOfAccounts(t *testing.T) {
	repository, _, _ := createRepository(t)

	// test that number of accounts is 0 after initialization
	if repository.GetNumberOfAccounts() != 0 {
		t.Errorf("expected 0, got %v", repository.GetNumberOfAccounts())
	}

	// test that number of accounts is set correctly
	repository.SetNumberOfAccounts(100)
	if repository.GetNumberOfAccounts() != 100 {
		t.Errorf("expected 100, got %v", repository.GetNumberOfAccounts())
	}
}

// createRepository creates a new repository instance with mocked dependencies.
func createRepository(t *testing.T) (*Repository, *rpc.MockRpc, *db.MockDatabase) {
	ctrl := gomock.NewController(t)
	mockRpc := rpc.NewMockRpc(ctrl)
	mockDb := db.NewMockDatabase(ctrl)
	repository := NewRepository(10_000, mockRpc, mockDb)
	return repository, mockRpc, mockDb
}

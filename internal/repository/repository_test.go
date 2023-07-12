package repository

import (
	"ftm-explorer/internal/repository/rpc"
	"ftm-explorer/internal/types"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

// Test that repository contains no blocks after initialization.
func TestRepository_NoObservedBlocks(t *testing.T) {
	repository, _ := createRepository(t)

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
func TestRepository_GetBlockByNumber(t *testing.T) {
	repository, mockRpc := createRepository(t)

	// block should be returned from rpc because it is not in the buffer
	block := types.Block{Number: 100}

	// expect rpc.BlockByNumber to be called with block number 100
	mockRpc.EXPECT().BlockByNumber(gomock.Eq(uint64(100))).Return(&block, nil)
	returnedBlock, err := repository.GetBlockByNumber(100)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if returnedBlock.Number != block.Number {
		t.Errorf("expected %v, got %v", block.Number, returnedBlock.Number)
	}

	// update latest observed block
	repository.UpdateLatestObservedBlock(&block)

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
	repository, mockRpc := createRepository(t)

	// transaction should be returned from rpc
	trx := types.Transaction{Hash: common.HexToHash("0x123")}
	mockRpc.EXPECT().TransactionByHash(gomock.Eq(trx.Hash)).Return(&trx, nil)
	returnedTrx, err := repository.GetTransactionByHash(trx.Hash)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if returnedTrx.Hash != trx.Hash {
		t.Errorf("expected %v, got %v", trx.Hash, returnedTrx.Hash)
	}
}

// createRepository creates a new repository instance with mocked dependencies.
func createRepository(t *testing.T) (*Repository, *rpc.MockRpc) {
	ctrl := gomock.NewController(t)
	mockRpc := rpc.NewMockRpc(ctrl)
	repository := NewRepository(mockRpc)
	return repository, mockRpc
}

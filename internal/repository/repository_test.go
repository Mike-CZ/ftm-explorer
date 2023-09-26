package repository

import (
	"ftm-explorer/internal/repository/db"
	"ftm-explorer/internal/repository/meta_fetcher"
	"ftm-explorer/internal/repository/rpc"
	"ftm-explorer/internal/types"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
)

// Test that repository contains no blocks after initialization.
func TestRepository_NoObservedBlocks(t *testing.T) {
	repository, _, _, _ := createRepository(t)

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
	repository, mockRpc, mockDb, _ := createRepository(t)

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
	repository, mockRpc, _, _ := createRepository(t)

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
	repository, mockRpc, _, _ := createRepository(t)

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
	repository, _, _, _ := createRepository(t)

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

// Test that the disk size per 100M transactions is set and returned correctly.
func TestRepository_GetAndSetDiskSizePer100MTxs(t *testing.T) {
	repository, _, _, _ := createRepository(t)

	// test that disk size per 100M transactions is 0 after initialization
	if repository.GetDiskSizePer100MTxs() != 0 {
		t.Errorf("expected 0, got %v", repository.GetDiskSizePer100MTxs())
	}

	// test that disk size per 100M transactions is set correctly
	repository.SetDiskSizePer100MTxs(289)
	if repository.GetDiskSizePer100MTxs() != 289 {
		t.Errorf("expected 289, got %v", repository.GetDiskSizePer100MTxs())
	}
}

// Test that repository transaction count is returned correctly.
func TestRepository_GetTrxCount(t *testing.T) {
	repository, _, mockDb, _ := createRepository(t)

	// transaction should be returned from rpc
	count := uint64(100)
	mockDb.EXPECT().TrxCount(gomock.Any()).Return(count, nil)
	returnedCount, err := repository.GetTrxCount()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if returnedCount != count {
		t.Errorf("expected %v, got %v", count, returnedCount)
	}
}

// Test that repository transaction count is called correctly.
func TestRepository_IncrementTrxCount(t *testing.T) {
	repository, _, mockDb, _ := createRepository(t)

	// check that trx increment is called on database
	mockDb.EXPECT().IncrementTrxCount(gomock.Any(), gomock.Eq(uint(10))).Return(nil)
	err := repository.IncrementTrxCount(10)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Test that repository returns number of validators.
func TestRepository_GetNumberOfValidators(t *testing.T) {
	repository, mockRpc, _, _ := createRepository(t)

	// check that number of validators method is called on rpc
	mockRpc.EXPECT().NumberOfValidators(gomock.Any()).Return(uint64(50), nil)
	count, err := repository.GetNumberOfValidators()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 50 {
		t.Errorf("expected 50, got %v", count)
	}
}

// Test that repository fetches number of accounts.
func TestRepository_FetchNumberOfAccounts(t *testing.T) {
	repository, _, _, mockFetcher := createRepository(t)

	// check that number of accounts method is called on meta fetcher
	mockFetcher.EXPECT().NumberOfAccounts().Return(uint64(11), nil)
	count, err := repository.FetchNumberOfAccounts()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if count != 11 {
		t.Errorf("expected 11, got %v", count)
	}
}

// Test that repository fetches disk size per 100M transactions.
func TestRepository_FetchDiskSizePer100MTxs(t *testing.T) {
	repository, _, _, mockFetcher := createRepository(t)

	// check that number of accounts method is called on meta fetcher
	mockFetcher.EXPECT().DiskSizePer100MTxs().Return(uint64(72799695667), nil)
	size, err := repository.FetchDiskSizePer100MTxs()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if size != 72799695667 {
		t.Errorf("expected 72799695667, got %v", size)
	}
}

// Test that repository fetches time to finality.
func TestRepository_FetchTimeToFinality(t *testing.T) {
	repository, _, _, mockFetcher := createRepository(t)

	// check that time to finality method is called on meta fetcher
	mockFetcher.EXPECT().TimeToFinality().Return(11.5, nil)
	time, err := repository.FetchTimeToFinality()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if time != 11.5 {
		t.Errorf("expected 11.5, got %v", time)
	}
}

// Test that signed transaction is sent.
func TestRepository_SendSignedTransaction(t *testing.T) {
	repository, mockRpc, _, _ := createRepository(t)

	// check that signed transaction is sent
	mockRpc.EXPECT().SendSignedTransaction(gomock.Any(), gomock.Any()).Return(nil)
	err := repository.SendSignedTransaction(&eth.Transaction{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Test that pending nonce is returned.
func TestRepository_PendingNonceAt(t *testing.T) {
	repository, mockRpc, _, _ := createRepository(t)

	// check that pending nonce is returned
	mockRpc.EXPECT().PendingNonceAt(gomock.Any(), gomock.Any()).Return(uint64(10), nil)
	nonce, err := repository.PendingNonceAt(common.Address{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if nonce != 10 {
		t.Errorf("expected 10, got %v", nonce)
	}
}

// Test that suggested gas price is returned.
func TestRepository_SuggestGasPrice(t *testing.T) {
	repository, mockRpc, _, _ := createRepository(t)

	// check that pending nonce is returned
	mockRpc.EXPECT().SuggestGasPrice(gomock.Any()).Return(big.NewInt(10), nil)
	gasPrice, err := repository.SuggestGasPrice()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if gasPrice.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("expected 10, got %v", gasPrice)
	}
}

// Test that network id is returned.
func TestRepository_NetworkID(t *testing.T) {
	repository, mockRpc, _, _ := createRepository(t)

	// check that network id is returned
	mockRpc.EXPECT().NetworkID(gomock.Any()).Return(big.NewInt(10), nil)
	networkID, err := repository.NetworkID()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if networkID.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("expected 10, got %v", networkID)
	}
}

// Test that time to block is calculated correctly.
func TestRepository_TimeToBlock(t *testing.T) {
	repository, _, _, _ := createRepository(t)

	// add 100 blocks to the buffer with 10 seconds between each block
	for i := 1; i <= 100; i++ {
		repository.blkBuffer.Add(
			&types.Block{
				Number:    hexutil.Uint64(i),
				Timestamp: hexutil.Uint64(i * 10)},
		)
	}

	// test that time to block is 7.5 seconds
	timeToBlock := repository.GetTimeToBlock()
	if timeToBlock != 10 {
		t.Errorf("expected 10, got %v", timeToBlock)
	}
}

// createRepository creates a new repository instance with mocked dependencies.
func createRepository(t *testing.T) (*Repository, *rpc.MockRpc, *db.MockDatabase, *meta_fetcher.MockMetaFetcher) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRpc := rpc.NewMockRpc(ctrl)
	mockDb := db.NewMockDatabase(ctrl)
	mockFetcher := meta_fetcher.NewMockMetaFetcher(ctrl)
	repository := NewRepository(10_000, mockRpc, mockDb, mockFetcher)
	return repository, mockRpc, mockDb, mockFetcher
}

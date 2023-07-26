package svc

import (
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
)

// Test block scanner run
func TestBlockScanner_Run(t *testing.T) {
	// initialize stubs
	ctrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepository(ctrl)
	mockLogger := logger.NewMockLogger()

	// create a channel for new headers, which will be used by the scanner
	heads := make(chan *eth.Header)
	mockRepository.EXPECT().GetNewHeadersChannel().Return(heads)

	// send heads into channel and expect them to be received by the scanner
	go func() {
		for i := 0; i <= 10; i++ {
			mockRepository.EXPECT().GetBlockByNumber(gomock.Eq(uint64(i))).Return(&types.Block{Number: hexutil.Uint64(i)}, nil)
			heads <- &eth.Header{Number: big.NewInt(int64(i))}
		}
	}()

	// start scanner
	scanner := newBlockScanner(&Manager{repo: mockRepository, log: mockLogger})
	scanner.start()
	defer scanner.close()

	// get output channel with scanned blocks
	scannedBlocks := scanner.scannedBlocks()

	// receive scanned blocks and check their numbers
	for i := 0; i <= 10; i++ {
		block := <-scannedBlocks
		if block.Number != hexutil.Uint64(i) {
			t.Errorf("expected block number %d, got %d", i, block.Number)
		}
	}
}

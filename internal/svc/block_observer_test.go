package svc

import (
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
)

func TestBlockObserver_Run(t *testing.T) {
	// initialize stubs
	ctrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepository(ctrl)
	mockLogger := logger.NewMockLogger()

	// create a channel, which will be used by the observer
	blocks := make(chan *types.Block)

	// start observer
	observer := newBlockObserver(&Manager{repo: mockRepository, log: mockLogger}, blocks)
	observer.start()
	defer observer.close()

	// validate, that observed blocks are forwarded to the repository
	for i := 0; i <= 10; i++ {
		blk := &types.Block{
			Number: hexutil.Uint64(i),
			Transactions: []common.Hash{
				common.HexToHash("0xabcd"),
				common.HexToHash("0x1234"),
			},
		}
		// expect the update of the latest observed block
		mockRepository.EXPECT().UpdateLatestObservedBlock(gomock.Eq(blk))

		// expect the update of the transactions count
		mockRepository.EXPECT().IncrementTrxCount(gomock.Eq(uint(2)))

		// send block to the observer
		blocks <- blk
	}
}

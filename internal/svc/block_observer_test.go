package svc

import (
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"testing"
	"time"

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
	observer := newBlockObserver(&Manager{repo: mockRepository, log: mockLogger, cfg: &config.Config{}}, blocks)
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

		// expect not to be idle
		mockRepository.EXPECT().IsIdle().Return(false)

		// expect the update of the transactions count
		mockRepository.EXPECT().IncrementTrxCount(gomock.Eq(uint(2)))

		// send block to the observer
		blocks <- blk
	}
}

// TestBlockObserver_IdleIsSet tests that the idle flag is set when the observer is idle.
func TestBlockObserver_IdleIsSet(t *testing.T) {
	// initialize stubs
	ctrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepository(ctrl)
	mockLogger := logger.NewMockLogger()

	// create a channel, which will be used by the observer
	blocks := make(chan *types.Block)

	// start observer
	observer := newBlockObserver(&Manager{repo: mockRepository, log: mockLogger, cfg: &config.Config{}}, blocks)

	// set timeout to 1 second
	observer.timeOutDuration = 1 * time.Second

	observer.start()
	defer observer.close()

	// expect that the idle flag is set
	mockRepository.EXPECT().IsIdle().Return(false)
	mockRepository.EXPECT().SetIsIdle(gomock.Eq(true))

	// wait for 1.5 second
	time.Sleep(1500 * time.Millisecond)
}

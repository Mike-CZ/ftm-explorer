package svc

import (
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"testing"

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
	observer := NewBlockObserver(blocks, mockRepository, mockLogger)
	observer.Start()
	defer observer.Stop()

	// validate, that observed blocks are forwarded to the repository
	for i := 0; i <= 10; i++ {
		blk := &types.Block{Number: hexutil.Uint64(i)}
		mockRepository.EXPECT().UpdateLatestObservedBlock(gomock.Eq(blk))
		blocks <- blk
	}
}

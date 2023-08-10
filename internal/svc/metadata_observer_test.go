package svc

import (
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

// Test metadata observer run
func TestMetadataObserver_Run(t *testing.T) {
	// initialize stubs
	ctrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepository(ctrl)
	mockLogger := logger.NewMockLogger()

	// start observer
	tickDuration := 100 * time.Millisecond
	observer := newMetadataObserver(&Manager{repo: mockRepository, log: mockLogger})
	observer.tickDuration = tickDuration
	observer.start()
	defer observer.close()

	// expect call to fetcher and repository
	numberOfAccounts := uint64(100)
	diskSizePer100MTxs := uint64(72799695667)
	mockRepository.EXPECT().FetchNumberOfAccounts().Return(numberOfAccounts, nil)
	mockRepository.EXPECT().SetNumberOfAccounts(gomock.Eq(numberOfAccounts))
	mockRepository.EXPECT().FetchDiskSizePer100MTxs().Return(diskSizePer100MTxs, nil)
	mockRepository.EXPECT().SetDiskSizePer100MTxs(gomock.Eq(diskSizePer100MTxs))

	// wait for ticker, add some extra time to make sure the ticker has ticked
	time.Sleep(tickDuration + tickDuration/2)
}

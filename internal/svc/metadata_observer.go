package svc

import (
	"ftm-explorer/internal/repository/meta_fetcher"
	"time"
)

// kMetadataTickDuration represents the frequency of the metadata observer progress.
const kMetadataTickDuration = 5 * time.Second

// metadataObserver represents the blockchain metadata observer.
type metadataObserver struct {
	service
	mf           meta_fetcher.IMetaFetcher
	sigClose     chan struct{}
	tickDuration time.Duration
}

// newMetadataObserver returns a new metadata observer.
func newMetadataObserver(mgr *Manager, mf meta_fetcher.IMetaFetcher) *metadataObserver {
	return &metadataObserver{
		service: service{
			mgr:  mgr,
			repo: mgr.repo,
			log:  mgr.log.ModuleLogger("metadata_observer"),
		},
		mf:           mf,
		sigClose:     make(chan struct{}, 1),
		tickDuration: kMetadataTickDuration,
	}
}

// start starts the metadata observer.
func (mo *metadataObserver) start() {
	mo.mgr.started(mo)
	go mo.execute()
}

// close stops the metadata observer.
func (mo *metadataObserver) close() {
	mo.sigClose <- struct{}{}
	mo.mgr.finished(mo)
}

// name returns the name of the metadata observer.
func (mo *metadataObserver) name() string {
	return "metadata_observer"
}

// execute executes the metadata observer.
func (mo *metadataObserver) execute() {
	ticker := time.NewTicker(mo.tickDuration)
	defer ticker.Stop()

	lastNumberOfAccounts := uint64(0)
	for {
		select {
		case <-mo.sigClose:
			return
		case <-ticker.C:
			numberOfAccounts, err := mo.mf.NumberOfAccounts()
			if err != nil {
				mo.log.Errorf("failed to get number of accounts: %v", err)
				continue
			}
			if numberOfAccounts != lastNumberOfAccounts {
				mo.log.Noticef("number of accounts: %d", numberOfAccounts)
				mo.repo.SetNumberOfAccounts(numberOfAccounts)
				lastNumberOfAccounts = numberOfAccounts
			}
		}
	}
}

package svc

import (
	"ftm-explorer/internal/types"
	"time"
)

// kMetadataTickDuration represents the frequency of the metadata observer progress.
const kMetadataTickDuration = 5 * time.Second

// metadataObserver represents the blockchain metadata observer.
type metadataObserver struct {
	service
	sigClose     chan struct{}
	tickDuration time.Duration

	lastTtfTime uint64
}

// newMetadataObserver returns a new metadata observer.
func newMetadataObserver(mgr *Manager) *metadataObserver {
	return &metadataObserver{
		service: service{
			mgr:  mgr,
			repo: mgr.repo,
			log:  mgr.log.ModuleLogger("metadata_observer"),
		},
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
	latestDiskSizePer100MTxs := uint64(0)
	latestDiskSizePrunedPer100MTxs := uint64(0)
	latestIsIdleStatus := false

	for {
		select {
		case <-mo.sigClose:
			return
		case <-ticker.C:
			// fetch and update number of accounts
			numberOfAccounts, err := mo.repo.FetchNumberOfAccounts()
			if err != nil {
				mo.log.Errorf("failed to get number of accounts: %v", err)
			} else if numberOfAccounts != lastNumberOfAccounts {
				mo.log.Noticef("number of accounts: %d", numberOfAccounts)
				mo.repo.SetNumberOfAccounts(numberOfAccounts)
				lastNumberOfAccounts = numberOfAccounts
			}

			// fetch and update disk size per 100M txs
			diskSizePer100MTxs, err := mo.repo.FetchDiskSizePer100MTxs()
			if err != nil {
				mo.log.Errorf("failed to get disk size per 100M txs: %v", err)
			} else if diskSizePer100MTxs != latestDiskSizePer100MTxs {
				mo.log.Noticef("disk size per 100M txs: %d", diskSizePer100MTxs)
				mo.repo.SetDiskSizePer100MTxs(diskSizePer100MTxs)
				latestDiskSizePer100MTxs = diskSizePer100MTxs
			}

			// fetch and update disk size pruned per 100M txs
			diskSizePrunedPer100MTxs, err := mo.repo.FetchDiskSizePrunedPer100MTxs()
			if err != nil {
				mo.log.Errorf("failed to get disk size pruned per 100M txs: %v", err)
			} else if diskSizePrunedPer100MTxs != latestDiskSizePrunedPer100MTxs {
				mo.log.Noticef("disk size pruned per 100M txs: %d", diskSizePrunedPer100MTxs)
				mo.repo.SetDiskSizePrunedPer100MTxs(diskSizePrunedPer100MTxs)
				latestDiskSizePrunedPer100MTxs = diskSizePrunedPer100MTxs
			}

			// fetch and update time to finality
			ttf, err := mo.repo.FetchTimeToFinality()
			if err != nil {
				mo.log.Errorf("failed to get time to finality: %v", err)
			} else {
				// add time to finality
				currentTime := time.Now().Unix()
				err := mo.repo.AddTimeToFinality(&types.Ttf{
					Timestamp: currentTime,
					Value:     ttf,
				})
				if err != nil {
					mo.log.Errorf("failed to add time to finality: %v", err)
				} else {
					resolution := types.AggResolutionSeconds
					seconds := uint64(resolution.ToDuration())
					aggTime := (uint64(currentTime) / seconds) * seconds
					if aggTime > mo.lastTtfTime {
						// calculate time to finality aggregations, get 60 ticks
						ttfAgg, err := mo.repo.GetTtfAvgAggByTimestamp(resolution, 60, aggTime)
						if err != nil {
							mo.log.Errorf("failed to get time to finality aggregation: %v", err)
						} else {
							mo.repo.SetTimeToFinalityPer10Secs(ttfAgg)
							mo.lastTtfTime = aggTime
						}
					}
					mo.log.Noticef("time to finality: %f", ttf)
				}
			}

			isIdleStatus, err := mo.repo.FetchIsIdleStatus()
			if err != nil {
				mo.log.Errorf("failed to get is idle status: %v", err)
			} else if isIdleStatus != latestIsIdleStatus {
				mo.log.Noticef("is idle status: %v", isIdleStatus)
				mo.repo.SetIsIdleOverride(isIdleStatus)
				latestIsIdleStatus = isIdleStatus
			}
		}
	}
}

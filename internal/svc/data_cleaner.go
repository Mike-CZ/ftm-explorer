package svc

import "time"

// kCleanTickDuration represents the frequency of the data cleaner default progress.
const kCleanTickDuration = 5 * time.Second

// kCleanMaxTtfCount represents the maximum number of ttfs to be persisted.
const kCleanMaxTtfCount = 1_000

// dataCleaner represents a cleaner of blockchain data.
type dataCleaner struct {
	service
	sigClose chan struct{}
}

// newDataCleaner creates a new data cleaner.
func newDataCleaner(mgr *Manager) *dataCleaner {
	return &dataCleaner{
		service: service{
			mgr:  mgr,
			repo: mgr.repo,
			log:  mgr.log.ModuleLogger("data_cleaner"),
		},
		sigClose: make(chan struct{}, 1),
	}
}

// start starts the data cleaner.
func (dc *dataCleaner) start() {
	dc.mgr.started(dc)
	go dc.execute()
}

// close stops the block scanner.
func (dc *dataCleaner) close() {
	dc.sigClose <- struct{}{}
	dc.mgr.finished(dc)
}

// name returns the name of the data cleaner.
func (dc *dataCleaner) name() string {
	return "data_cleaner"
}

// execute executes the data cleaner.
func (dc *dataCleaner) execute() {
	ticker := time.NewTicker(kCleanTickDuration)
	defer ticker.Stop()

	for {
		select {
		case <-dc.sigClose:
			return
		case <-ticker.C:
			dc.cleanTransactions()
			dc.cleanTtf()
		}
	}
}

// cleanTransactions cleans the transactions collection.
func (dc *dataCleaner) cleanTransactions() {
	// do not clean if the explorer is not persisted
	if !dc.mgr.cfg.Explorer.IsPersisted {
		return
	}
	if err := dc.repo.ShrinkTransactions(int64(dc.mgr.cfg.Explorer.MaxTxsCount)); err != nil {
		dc.log.Errorf("failed to shrink transactions: %s", err.Error())
	}
}

// cleanTransactions cleans the transactions collection.
func (dc *dataCleaner) cleanTtf() {
	if err := dc.repo.ShrinkTtf(kCleanMaxTtfCount); err != nil {
		dc.log.Errorf("failed to shrink ttf: %s", err.Error())
	}
}

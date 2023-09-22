package svc

import (
	"ftm-explorer/internal/types"
)

// blockObserver represents an observer of blockchain blocks.
type blockObserver struct {
	service
	inBlocks <-chan *types.Block
	sigClose chan struct{}

	// lastAggTime is the last time the aggregator was run.
	lastAggTime uint64
}

// newBlockObserver creates a new block observer.
// It observes new blocks which are sent to the channel. It then processes them.
func newBlockObserver(mgr *Manager, inBlocks <-chan *types.Block) *blockObserver {
	return &blockObserver{
		service: service{
			mgr:  mgr,
			repo: mgr.repo,
			log:  mgr.log.ModuleLogger("block_observer"),
		},
		inBlocks: inBlocks,
		sigClose: make(chan struct{}, 1),
	}
}

// start starts the block observer.
func (bs *blockObserver) start() {
	bs.mgr.started(bs)
	go bs.execute()
}

// close stops the block observer.
func (bs *blockObserver) close() {
	bs.sigClose <- struct{}{}
	bs.mgr.finished(bs)
}

// name returns the name of the block observer.
func (bs *blockObserver) name() string {
	return "block_observer"
}

// execute executes the block observer.
func (bs *blockObserver) execute() {
	for {
		select {
		case <-bs.sigClose:
			return
		case block, ok := <-bs.inBlocks:
			if !ok {
				bs.log.Notice("input blocks channel closed. stopping block observer")
				return
			}
			bs.processBlock(block)
		}
	}
}

// processBlock processes a block.
func (bs *blockObserver) processBlock(block *types.Block) {
	bs.log.Noticef("block observer processing block %d", block.Number)

	// update latest observed block
	if err := bs.repo.UpdateLatestObservedBlock(block); err != nil {
		bs.log.Errorf("error updating latest observed block: %v", err)
		return
	}

	// increment transaction count
	if err := bs.repo.IncrementTrxCount(uint(len(block.Transactions))); err != nil {
		bs.log.Errorf("error incrementing transaction count: %v", err)
		return
	}

	// get block time rounded to nearest seconds
	// in case it is slowing down the whole process too much, we can move it to a separate goroutine
	resolution := types.AggResolutionSeconds
	seconds := uint64(resolution.ToDuration())
	aggTime := (uint64(block.Timestamp) / seconds) * seconds

	// if the time is different from the last time the aggregator was run, run the aggregations
	// we also get last 60 ticks (10 seconds each) to be used by the chart
	// also wait for the latest block to be at least 1 second older than the aggregation time
	// so that we don't miss any blocks
	if aggTime != bs.lastAggTime && uint64(block.Timestamp) > aggTime+1 {
		bs.lastAggTime = aggTime
		txAggs, err := bs.repo.GetTrxCountAggByTimestamp(resolution, 60, &aggTime)
		if err != nil {
			bs.log.Errorf("error getting transaction count aggregation by timestamp: %v", err)
			return
		}
		gasUsedAggs, err := bs.repo.GetGasUsedAggByTimestamp(resolution, 60, &aggTime)
		if err != nil {
			bs.log.Errorf("error getting gas used aggregation by timestamp: %v", err)
			return
		}
		bs.repo.SetTxCountPer10Secs(txAggs)
		bs.repo.SetGasUsedPer10Secs(gasUsedAggs)
		bs.log.Notice("aggregation data updated successfully")
	}
}

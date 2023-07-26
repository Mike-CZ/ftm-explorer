package svc

import (
	"ftm-explorer/internal/types"
)

// blockObserver represents an observer of blockchain blocks.
type blockObserver struct {
	service
	inBlocks <-chan *types.Block
	sigClose chan struct{}
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
}

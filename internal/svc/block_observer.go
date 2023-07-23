package svc

import (
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"sync"
)

// BlockObserver represents an observer of blockchain blocks.
type BlockObserver struct {
	repo     repository.IRepository
	log      logger.ILogger
	inBlocks <-chan *types.Block
	sigClose chan struct{}
	wg       sync.WaitGroup
}

// NewBlockObserver creates a new block observer.
// It observes new blocks which are sent to the channel. It then processes them.
func NewBlockObserver(inBlocks <-chan *types.Block, repo repository.IRepository, log logger.ILogger) *BlockObserver {
	return &BlockObserver{
		repo:     repo,
		log:      log.ModuleLogger("block_observer"),
		inBlocks: inBlocks,
		sigClose: make(chan struct{}, 1),
	}
}

// Start starts the block observer.
func (bs *BlockObserver) Start() {
	bs.wg.Add(1)
	go bs.execute()
}

// Stop stops the block observer.
func (bs *BlockObserver) Stop() {
	bs.sigClose <- struct{}{}
	bs.wg.Wait()
}

// execute executes the block observer.
func (bs *BlockObserver) execute() {
	bs.log.Notice("block observer started")
	defer bs.wg.Done()

	for {
		select {
		case <-bs.sigClose:
			bs.log.Notice("block observer stopped")
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
func (bs *BlockObserver) processBlock(block *types.Block) {
	bs.log.Noticef("block observer processing block %d", block.Number)

	// update latest observed block
	if err := bs.repo.UpdateLatestObservedBlock(block); err != nil {
		bs.log.Errorf("error updating latest observed block: %v", err)
		return
	}
}

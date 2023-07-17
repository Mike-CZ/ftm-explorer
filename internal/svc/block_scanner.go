package svc

import (
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"sync"
	"time"
)

// kOutBlockBufferCapacity represents the capacity of the found blocks channel.
const kOutBlockBufferCapacity = 10_000

// kScanTickDuration represents the frequency of the scanner default progress.
const kScanTickDuration = 5 * time.Millisecond

// BlockScanner represents a scanner of blockchain blocks.
// It scans the blockchain for new blocks and sends them to the channel.
type BlockScanner struct {
	repo      repository.IRepository
	log       logger.ILogger
	outBlocks chan *types.Block
	sigClose  chan struct{}
	wg        sync.WaitGroup
}

// NewBlockScanner creates a new block scanner.
func NewBlockScanner(repo repository.IRepository, log logger.ILogger) *BlockScanner {
	return &BlockScanner{
		repo:      repo,
		log:       log.ModuleLogger("block_scanner"),
		outBlocks: make(chan *types.Block, kOutBlockBufferCapacity),
		sigClose:  make(chan struct{}, 1),
	}
}

// ScannedBlocks returns a channel containing scanned blocks.
func (bs *BlockScanner) ScannedBlocks() <-chan *types.Block {
	return bs.outBlocks
}

// Start starts the block scanner.
func (bs *BlockScanner) Start() {
	bs.wg.Add(1)
	go bs.execute()
}

// Stop stops the block scanner.
func (bs *BlockScanner) Stop() {
	bs.sigClose <- struct{}{}
	bs.wg.Wait()
}

// execute executes the block scanner.
func (bs *BlockScanner) execute() {
	bs.log.Notice("block scanner started")
	defer bs.wg.Done()

	// get channel with new headers
	heads := bs.repo.GetNewHeadersChannel()

	// start ticker
	ticker := time.NewTicker(kScanTickDuration)
	defer ticker.Stop()

	var targetBlock *uint64
	var nextBlock *uint64
	for {
		select {
		// we should stop
		case <-bs.sigClose:
			bs.log.Notice("block scanner stopped")
			return
		// we have a new target block
		case head, ok := <-heads:
			if !ok {
				bs.log.Notice("new headers channel closed. stopping block scanner")
				return
			}
			// initialize target block if it is not initialized
			if targetBlock == nil {
				targetBlock = new(uint64)
			}
			*targetBlock = head.Number.Uint64()
			bs.log.Debugf("block scanner target block set to %d", targetBlock)
			// if we have no next block, set it to the target block
			if nextBlock == nil {
				nextBlock = new(uint64)
				*nextBlock = *targetBlock
			}
		// scan new blocks
		case <-ticker.C:
			if targetBlock != nil && nextBlock != nil && *targetBlock >= *nextBlock {
				block, err := bs.repo.GetBlockByNumber(*nextBlock)
				if err != nil {
					bs.log.Warningf("block scanner can not proceed; %v", err)
					continue
				}
				bs.outBlocks <- block
				*nextBlock++
			}
		}
	}
}

package svc

import (
	"ftm-explorer/internal/types"
	"time"
)

// kOutBlockBufferCapacity represents the capacity of the found blocks channel.
const kOutBlockBufferCapacity = 10_000

// kScanTickDuration represents the frequency of the scanner default progress.
const kScanTickDuration = 5 * time.Millisecond

// blockScanner represents a scanner of blockchain blocks.
// It scans the blockchain for new blocks and sends them to the channel.
type blockScanner struct {
	service
	outBlocks chan *types.Block
	sigClose  chan struct{}
}

// newBlockScanner creates a new block scanner.
func newBlockScanner(mgr *Manager) *blockScanner {
	return &blockScanner{
		service: service{
			mgr:  mgr,
			repo: mgr.repo,
			log:  mgr.log.ModuleLogger("block_scanner"),
		},
		outBlocks: make(chan *types.Block, kOutBlockBufferCapacity),
		sigClose:  make(chan struct{}, 1),
	}
}

// scannedBlocks returns a channel containing scanned blocks.
func (bs *blockScanner) scannedBlocks() <-chan *types.Block {
	return bs.outBlocks
}

// start starts the block scanner.
func (bs *blockScanner) start() {
	bs.mgr.started(bs)
	go bs.execute()
}

// close stops the block scanner.
func (bs *blockScanner) close() {
	bs.sigClose <- struct{}{}
	bs.mgr.finished(bs)
}

func (bs *blockScanner) name() string {
	return "block_scanner"
}

// execute executes the block scanner.
func (bs *blockScanner) execute() {
	// get channel with new headers
	heads := bs.repo.GetNewHeadersChannel()

	// start ticker
	ticker := time.NewTicker(kScanTickDuration)
	defer ticker.Stop()

	var targetBlock *uint64
	var nextBlock *uint64
	for {
		select {
		// we should close
		case <-bs.sigClose:
			return
		// we have a new target block
		case head, ok := <-heads:
			if !ok {
				bs.log.Notice("new headers channel closed. stopping block scanner")
				bs.mgr.finished(bs)
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
				if block == nil {
					bs.log.Warningf("block scanner can not proceed; block %d not found", *nextBlock)
					continue
				}
				bs.outBlocks <- block
				*nextBlock++
			}
		}
	}
}

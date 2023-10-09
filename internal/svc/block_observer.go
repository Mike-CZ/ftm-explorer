package svc

import (
	db_types "ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"
	"sync"
)

// blockObserver represents an observer of blockchain blocks.
type blockObserver struct {
	service
	inBlocks <-chan *types.Block
	sigClose chan struct{}

	// lastAggTime is the last time the aggregator was run.
	lastAggTime uint64

	// aggMtx is a mutex used to synchronize access to the aggregator
	aggMtx sync.Mutex
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
	wg := sync.WaitGroup{}

	for {
		select {
		case <-bs.sigClose:
			// wait for all goroutines to finish, then return
			wg.Wait()
			return
		case block, ok := <-bs.inBlocks:
			if !ok {
				bs.log.Notice("input blocks channel closed. stopping block observer")
				return
			}
			wg.Add(1)
			go bs.processBlock(block, &wg)
		}
	}
}

// processBlock processes a block.
func (bs *blockObserver) processBlock(block *types.Block, wg *sync.WaitGroup) {
	defer wg.Done()

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
		// lock the aggregator, so that only one goroutine can run it at a time
		updateAggs := bs.aggMtx.TryLock()
		if updateAggs {
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

			// unlock the aggregator
			bs.aggMtx.Unlock()
		}
	}

	// store transactions
	if bs.mgr.cfg.Explorer.IsPersisted {
		bs.storeTransactions(block)
	}
}

// storeTransactions stores transactions in the database.
func (bs *blockObserver) storeTransactions(block *types.Block) {
	var txs []db_types.Transaction

	if len(block.Transactions) == 0 {
		bs.log.Noticef("no transactions to store in block %d", block.Number)
		return
	}

	for _, hash := range block.Transactions {
		// get transaction
		tx, err := bs.repo.GetTransactionByHash(hash)
		if err != nil {
			bs.log.Errorf("error getting transaction %s: %v", hash, err)
			continue
		}
		if tx == nil {
			bs.log.Errorf("transaction %s not found", hash)
			continue
		}
		dbTx := db_types.Transaction{
			Hash:      tx.Hash,
			Timestamp: int64(block.Timestamp),
		}
		// append sender address
		dbTx.Addresses = append(dbTx.Addresses, tx.From)
		// append receiver address if it is not a contract creation
		if tx.To != nil {
			dbTx.Addresses = append(dbTx.Addresses, *tx.To)
		}
		// append contract address if it is a contract creation
		if tx.ContractAddress != nil {
			dbTx.Addresses = append(dbTx.Addresses, *tx.ContractAddress)
		}
		// append transaction to the list
		txs = append(txs, dbTx)
	}

	// store transactions
	if err := bs.repo.AddTransactions(txs); err != nil {
		bs.log.Criticalf("error storing transactions: %v", err)
		return
	}

	bs.log.Noticef("stored %d transactions for block %d", len(txs), block.Number)
}

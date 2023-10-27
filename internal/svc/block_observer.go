package svc

import (
	db_types "ftm-explorer/internal/repository/db/types"
	"ftm-explorer/internal/types"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// kObserverChainTimeOutDuration represents the timeout duration of the observer chain.
const kObserverChainTimeOutDuration = 5 * time.Second

// blockObserver represents an observer of blockchain blocks.
type blockObserver struct {
	service
	inBlocks <-chan *types.Block
	sigClose chan struct{}

	// lastAggTime is the last time the aggregator was run.
	lastAggTime uint64

	// lastBlkTime is the last time a block was processed.
	lastBlkTime uint64

	// timeOutDuration is the timeout duration of the observer chain.
	timeOutDuration time.Duration

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
		inBlocks:        inBlocks,
		sigClose:        make(chan struct{}, 1),
		timeOutDuration: kObserverChainTimeOutDuration,
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
	ticker := time.NewTicker(bs.timeOutDuration)
	defer ticker.Stop()

	for {
		select {
		case <-bs.sigClose:
			// wait for all goroutines to finish, then return
			wg.Wait()
			return
		case <-ticker.C:
			// if the last block time is older than the timeout duration, set the idle flag
			if !bs.repo.IsIdle() && uint64(time.Now().Unix())-bs.lastBlkTime >= uint64(bs.timeOutDuration.Seconds()) {
				bs.repo.SetIsIdle(true)
			}
		case block, ok := <-bs.inBlocks:
			if !ok {
				bs.log.Notice("input blocks channel closed. stopping block observer")
				return
			}
			bs.lastBlkTime = uint64(time.Now().Unix())
			// reset the idle flag
			if bs.repo.IsIdle() {
				bs.repo.SetIsIdle(false)
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
			// unlock the aggregator
			defer bs.aggMtx.Unlock()

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

	// store transactions
	if bs.mgr.cfg.Explorer.IsPersisted {
		bs.storeTransactions(block)
	}
}

// storeTransactions stores transactions in the database.
func (bs *blockObserver) storeTransactions(block *types.Block) {
	var txs []db_types.Transaction
	accounts := make(map[common.Address]bool)

	if len(block.Transactions) == 0 {
		bs.log.Noticef("no transactions to store in block %d", block.Number)
		return
	}

	for _, hash := range block.Transactions {
		txAccounts := make(map[common.Address]bool)

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
		txAccounts[tx.From] = true
		// append receiver address if it is not a contract creation
		if tx.To != nil {
			txAccounts[*tx.To] = true
		}
		// append contract address if it is a contract creation
		if tx.ContractAddress != nil {
			txAccounts[*tx.ContractAddress] = true
		}

		// handles events
		for _, log := range tx.Logs {
			if len(log.Topics) == 0 {
				continue
			}
			switch log.Topics[0].Hex() {
			// ERC20::Approval(address indexed owner, address indexed spender, uint256 value)
			case "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925":
				fallthrough
			// ERC20::Transfer(address indexed from, address indexed to, uint256 value)
			case "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef":
				if len(log.Topics) == 3 && len(log.Data) == 32 {
					from := common.BytesToAddress(log.Topics[1].Bytes())
					to := common.BytesToAddress(log.Topics[2].Bytes())
					txAccounts[from] = true
					txAccounts[to] = true
				}
			// UniswapPair::Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
			case "0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822":
				if len(log.Data) == 128 && len(log.Topics) == 3 {
					from := common.BytesToAddress(log.Topics[1].Bytes())
					to := common.BytesToAddress(log.Topics[2].Bytes())
					txAccounts[from] = true
					txAccounts[to] = true
				}
			}
		}

		// add addresses to the map
		for addr := range txAccounts {
			dbTx.Addresses = append(dbTx.Addresses, addr)
			accounts[addr] = true
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

	// store accounts
	var accountsList []common.Address
	for addr := range accounts {
		accountsList = append(accountsList, addr)
	}
	if err := bs.repo.AddAccounts(accountsList, int64(block.Timestamp)); err != nil {
		bs.log.Criticalf("error storing accounts: %v", err)
		return
	}
}

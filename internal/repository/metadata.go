package repository

import (
	"context"
	"ftm-explorer/internal/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// GetNumberOfAccounts returns the number of accounts in the blockchain.
func (r *Repository) GetNumberOfAccounts() uint64 {
	return r.numberOfAccounts
}

// SetNumberOfAccounts sets the number of accounts in the blockchain.
func (r *Repository) SetNumberOfAccounts(number uint64) {
	r.numberOfAccounts = number
}

// GetDiskSizePer100MTxs returns the disk size per 100M transactions.
func (r *Repository) GetDiskSizePer100MTxs() uint64 {
	return r.diskSizePer100MTxs
}

// SetDiskSizePer100MTxs sets the disk size per 100M transactions.
func (r *Repository) SetDiskSizePer100MTxs(number uint64) {
	r.diskSizePer100MTxs = number
}

// GetDiskSizePrunedPer100MTxs returns the disk size pruned per 100M transactions.
func (r *Repository) GetDiskSizePrunedPer100MTxs() uint64 {
	return r.diskSizePrunedPer100MTxs
}

// SetDiskSizePrunedPer100MTxs sets the disk size pruned per 100M transactions.
func (r *Repository) SetDiskSizePrunedPer100MTxs(number uint64) {
	r.diskSizePrunedPer100MTxs = number
}

// PendingNonceAt returns the nonce of the account at the given block.
func (r *Repository) PendingNonceAt(address common.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()

	return r.rpc.PendingNonceAt(ctx, address)
}

// SuggestGasPrice suggests a gas price.
func (r *Repository) SuggestGasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()

	return r.rpc.SuggestGasPrice(ctx)
}

// NetworkID returns the network ID.
func (r *Repository) NetworkID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kRpcTimeout)
	defer cancel()

	return r.rpc.NetworkID(ctx)
}

// GetTxCountPer10Secs returns transactions per 10 seconds.
func (r *Repository) GetTxCountPer10Secs() []types.HexUintTick {
	return r.txCountPer10Secs
}

// SetTxCountPer10Secs sets transactions per 10 seconds.
func (r *Repository) SetTxCountPer10Secs(data []types.HexUintTick) {
	cpy := make([]types.HexUintTick, len(data))
	copy(cpy, data)
	r.txCountPer10Secs = cpy
}

// GetGasUsedPer10Secs returns gas used per 10 seconds.
func (r *Repository) GetGasUsedPer10Secs() []types.HexUintTick {
	return r.gasUsedPer10Secs
}

// SetGasUsedPer10Secs sets gas used per 10 seconds.
func (r *Repository) SetGasUsedPer10Secs(data []types.HexUintTick) {
	cpy := make([]types.HexUintTick, len(data))
	copy(cpy, data)
	r.gasUsedPer10Secs = cpy
}

// AddTimeToFinality adds the given time to finality.
func (r *Repository) AddTimeToFinality(ttf *types.Ttf) error {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()

	return r.db.AddTimeToFinality(ctx, ttf)
}

// GetTtfAvgAggByTimestamp returns average aggregation of time to finality in given time range.
func (r *Repository) GetTtfAvgAggByTimestamp(resolution types.AggResolution, ticks uint, endTime uint64) ([]types.FloatTick, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()

	return r.db.TtfAvgAggByTimestamp(ctx, endTime, resolution.ToDuration(), ticks)
}

// GetTimeToFinality returns the time to finality in the blockchain.
func (r *Repository) GetTimeToFinality() float64 {
	// if there is no data, return 0
	if len(r.ttfPer10Secs) == 0 {
		return 0
	}

	// otherwise return last value rounded to 2 decimals
	return float64(int(r.ttfPer10Secs[len(r.ttfPer10Secs)-1].Value*100)) / 100
}

// GetTimeToFinalityPer10Secs returns time to finality per 10 seconds.
func (r *Repository) GetTimeToFinalityPer10Secs() []types.FloatTick {
	return r.ttfPer10Secs
}

// SetTimeToFinalityPer10Secs sets time to finality per 10 seconds.
func (r *Repository) SetTimeToFinalityPer10Secs(data []types.FloatTick) {
	cpy := make([]types.FloatTick, len(data))
	copy(cpy, data)
	r.ttfPer10Secs = cpy
}

// IsIdle returns isIdle.
func (r *Repository) IsIdle() bool {
	return r.isIdle
}

// SetIsIdle sets isIdle.
func (r *Repository) SetIsIdle(isIdle bool) {
	r.isIdle = isIdle
}

// IsIdleOverride returns isIdleOverride.
func (r *Repository) IsIdleOverride() bool {
	return r.isIdleOverride
}

// SetIsIdleOverride sets isIdleOverride.
func (r *Repository) SetIsIdleOverride(isIdleOverride bool) {
	r.isIdleOverride = isIdleOverride
}

// GetTimeToBlock returns the time to block.
func (r *Repository) GetTimeToBlock() float64 {
	// get last 100 blocks from buffer
	blocks := r.blkBuffer.GetLatest(100)
	// if there are no blocks or just 1 block, return 0
	if len(blocks) <= 1 {
		return 0
	}
	// calculate deltas between blocks
	totalDelta := 0
	for i := 1; i < len(blocks); i++ {
		totalDelta += int(blocks[i-1].Timestamp - blocks[i].Timestamp)
	}

	// calculate average delta
	avgDelta := float32(totalDelta) / float32(len(blocks)-1)

	// round to 2 decimals and return
	return float64(int(avgDelta*100)) / 100
}

// ShrinkTtf shrinks the time to finality collection. It will persist the given number of ttfs.
func (r *Repository) ShrinkTtf(count int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), kDbTimeout)
	defer cancel()

	return r.db.ShrinkTtf(ctx, count)
}

package resolvers

import (
	"ftm-explorer/internal/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// CurrentState resolves the current state of the blockchain.
type CurrentState struct {
	rs *RootResolver
}

// TtfTick represents a time to finality tick.
type TtfTick types.FloatTick

// State resolves the current state of the blockchain.
func (rs *RootResolver) State() CurrentState {
	return CurrentState{rs: rs}
}

// NumberOfAccounts returns the number of accounts.
func (rs *RootResolver) NumberOfAccounts() (int32, error) {
	if rs.isPersisted {
		num, err := rs.repository.GetNumberOfAccountsInDb()
		if err != nil {
			return 0, err
		}
		return int32(num), nil
	}
	return int32(rs.repository.GetNumberOfAccounts()), nil
}

// DiskSizePer100MTxs resolves the disk size per 100M transactions.
func (rs *RootResolver) DiskSizePer100MTxs() hexutil.Uint64 {
	return hexutil.Uint64(rs.repository.GetDiskSizePer100MTxs())
}

// DiskSizePrunedPer100MTxs resolves the disk pruned size per 100M transactions.
func (rs *RootResolver) DiskSizePrunedPer100MTxs() hexutil.Uint64 {
	return hexutil.Uint64(rs.repository.GetDiskSizePrunedPer100MTxs())
}

// NumberOfValidators resolves the number of validators in the blockchain.
func (rs *RootResolver) NumberOfValidators() int32 {
	number, err := rs.repository.GetNumberOfValidators()
	if err != nil {
		rs.log.Errorf("failed to get number of validators: %v", err)
		return 0
	}
	return int32(number)
}

// TtfTimestampAggregations resolves ttf timestamp aggregations.
func (rs *RootResolver) TtfTimestampAggregations() []TtfTick {
	result := rs.repository.GetTimeToFinalityPer10Secs()

	// convert result
	rv := make([]TtfTick, len(result))
	for i, t := range result {
		rv[i] = (TtfTick)(t)
	}

	return rv
}

// TimeToFinality resolves the time to finality.
func (rs *RootResolver) TimeToFinality() float64 {
	return rs.repository.GetTimeToFinality()
}

// TimeToBlock resolves the time to block.
func (rs *RootResolver) TimeToBlock() float64 {
	return rs.repository.GetTimeToBlock()
}

// IsIdle resolves if the chain is idle.
func (rs *RootResolver) IsIdle() bool {
	return rs.repository.IsIdleOverride() || rs.repository.IsIdle()
}

// CurrentBlockHeight resolves the current block height.
func (cs CurrentState) CurrentBlockHeight() (*hexutil.Uint64, error) {
	return cs.rs.CurrentBlockHeight()
}

// NumberOfAccounts resolves the number of accounts in the blockchain.
func (cs CurrentState) NumberOfAccounts() (int32, error) {
	return cs.rs.NumberOfAccounts()
}

// NumberOfTransactions resolves the number of transactions in the blockchain.
func (cs CurrentState) NumberOfTransactions() (hexutil.Uint64, error) {
	return cs.rs.NumberOfTransactions()
}

// NumberOfValidators resolves the number of validators in the blockchain.
func (cs CurrentState) NumberOfValidators() int32 {
	return cs.rs.NumberOfValidators()
}

// DiskSizePer100MTxs resolves the disk size per 100M transactions.
func (cs CurrentState) DiskSizePer100MTxs() hexutil.Uint64 {
	return cs.rs.DiskSizePer100MTxs()
}

// DiskSizePrunedPer100MTxs resolves the disk size pruned per 100M transactions.
func (cs CurrentState) DiskSizePrunedPer100MTxs() hexutil.Uint64 {
	return cs.rs.DiskSizePrunedPer100MTxs()
}

// TimeToFinality resolves the time to finality.
func (cs CurrentState) TimeToFinality() float64 {
	return cs.rs.TimeToFinality()
}

// TimeToBlock resolves the time to block.
func (cs CurrentState) TimeToBlock() float64 {
	return cs.rs.TimeToBlock()
}

// IsIdle resolves if the chain is idle.
func (cs CurrentState) IsIdle() bool {
	return cs.rs.IsIdle()
}

// Timestamp resolves tick timestamp.
func (t TtfTick) Timestamp() int32 {
	return int32(t.Time)
}

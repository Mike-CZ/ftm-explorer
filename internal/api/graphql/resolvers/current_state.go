package resolvers

import (
	"math"
	"math/rand"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// CurrentState resolves the current state of the blockchain.
type CurrentState struct {
	rs *RootResolver
}

// State resolves the current state of the blockchain.
func (rs *RootResolver) State() CurrentState {
	return CurrentState{rs: rs}
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

// DiskSizePer100MTxs resolves the disk size per 100M transactions.
func (rs *RootResolver) DiskSizePer100MTxs() hexutil.Uint64 {
	// ~ 54.5 GB per 100M transactions
	return 54_494_722_457
}

// TimeToFinality resolves the time to finality.
func (rs *RootResolver) TimeToFinality() float64 {
	min := 0.8
	max := 1.4
	randomValue := min + rand.Float64()*(max-min)
	roundedValue := math.Round(randomValue*100) / 100 // Round to 2 decimal places
	return roundedValue
}

// CurrentBlockHeight resolves the current block height.
func (cs CurrentState) CurrentBlockHeight() (*hexutil.Uint64, error) {
	return cs.rs.CurrentBlockHeight()
}

// NumberOfAccounts resolves the number of accounts in the blockchain.
func (cs CurrentState) NumberOfAccounts() int32 {
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

// TimeToFinality resolves the time to finality.
func (cs CurrentState) TimeToFinality() float64 {
	return cs.rs.TimeToFinality()
}

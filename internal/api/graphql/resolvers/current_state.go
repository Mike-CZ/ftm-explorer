package resolvers

import "github.com/ethereum/go-ethereum/common/hexutil"

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
	return 8
}

// DiskSizePer100MTxs resolves the disk size per 100M transactions.
func (rs *RootResolver) DiskSizePer100MTxs() hexutil.Uint64 {
	// ~ 54.5 GB per 100M transactions
	return 54_494_722_457
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

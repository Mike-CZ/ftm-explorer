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

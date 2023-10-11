package resolvers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Account represents resolvable blockchain account structure.
type Account struct {
	Address common.Address
	rs      *RootResolver
}

// Account represents resolvable blockchain account structure.
func (rs *RootResolver) Account(args *struct{ Address common.Address }) Account {
	return Account{Address: args.Address, rs: rs}
}

// Balance returns the balance of the account.
func (acc Account) Balance() (hexutil.Big, error) {
	val, err := acc.rs.repository.AccountBalance(acc.Address)
	if err != nil {
		return hexutil.Big{}, err
	}
	return *val, nil
}

// Transactions returns the transactions of the account.
func (acc Account) Transactions() ([]*Transaction, error) {
	// fetch 500 latest transactions
	txs, err := acc.rs.repository.GetLastTransactionsWhereAddress(acc.Address, 500)
	if err != nil {
		return nil, err
	}

	rv := make([]*Transaction, len(txs))
	for i, tx := range txs {
		t, err := acc.rs.repository.GetTransactionByHash(tx.Hash)
		if err != nil {
			return nil, err
		}
		rv[i] = (*Transaction)(t)
	}

	return rv, nil
}

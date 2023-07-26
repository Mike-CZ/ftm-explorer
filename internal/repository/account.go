package repository

// GetNumberOfAccounts returns the number of accounts in the blockchain.
func (r *Repository) GetNumberOfAccounts() uint64 {
	return r.numberOfAccounts
}

// SetNumberOfAccounts sets the number of accounts in the blockchain.
func (r *Repository) SetNumberOfAccounts(number uint64) {
	r.numberOfAccounts = number
}

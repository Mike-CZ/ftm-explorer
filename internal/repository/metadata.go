package repository

import "ftm-explorer/internal/types"

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

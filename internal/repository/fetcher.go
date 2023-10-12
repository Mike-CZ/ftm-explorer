package repository

// FetchNumberOfAccounts returns the number of accounts in the blockchain.
// This method will fetch data from remote host.
func (r *Repository) FetchNumberOfAccounts() (uint64, error) {
	return r.metaFetcher.NumberOfAccounts()
}

// FetchDiskSizePer100MTxs returns the disk size per 100M transactions.
// This method will fetch data from remote host.
func (r *Repository) FetchDiskSizePer100MTxs() (uint64, error) {
	return r.metaFetcher.DiskSizePer100MTxs()
}

// FetchDiskSizePrunedPer100MTxs returns the disk size pruned per 100M transactions.
// This method will fetch data from remote host.
func (r *Repository) FetchDiskSizePrunedPer100MTxs() (uint64, error) {
	return r.metaFetcher.DiskSizePrunedPer100MTxs()
}

// FetchTimeToFinality returns the time to finality in the blockchain.
// This method will fetch data from remote host.
func (r *Repository) FetchTimeToFinality() (float64, error) {
	return r.metaFetcher.TimeToFinality()
}

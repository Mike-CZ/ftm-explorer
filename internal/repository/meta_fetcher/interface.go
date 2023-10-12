package meta_fetcher

//go:generate mockgen -source=interface.go -destination=meta_fetcher_mock.go -package=meta_fetcher -mock_names=IMetaFetcher=MockMetaFetcher

// IMetaFetcher represents a meta fetcher. It is responsible for fetching the blockchain metadata.
type IMetaFetcher interface {
	// NumberOfAccounts returns the number of accounts in the blockchain.
	NumberOfAccounts() (uint64, error)
	// DiskSizePer100MTxs returns the disk size per 100M transactions.
	DiskSizePer100MTxs() (uint64, error)
	// DiskSizePrunedPer100MTxs returns the disk size pruned per 100M transactions.
	DiskSizePrunedPer100MTxs() (uint64, error)
	// TimeToFinality returns the time to finality in the blockchain.
	TimeToFinality() (float64, error)
}

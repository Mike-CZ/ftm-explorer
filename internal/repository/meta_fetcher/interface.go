package meta_fetcher

// IMetaFetcher represents a meta fetcher. It is responsible for fetching the blockchain metadata.
type IMetaFetcher interface {
	// NumberOfAccounts returns the number of accounts in the blockchain.
	NumberOfAccounts() (uint64, error)
}

package meta_fetcher

import (
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// MetaFetcher represents a meta fetcher. It is responsible for fetching the blockchain metadata.
type MetaFetcher struct {
	log logger.ILogger
	url string
}

// NewMetaFetcher returns a new meta fetcher.
func NewMetaFetcher(cfg *config.MetaFetcher, log logger.ILogger) *MetaFetcher {
	return &MetaFetcher{
		log: log.ModuleLogger("meta_fetcher"),
		url: cfg.Url,
	}
}

// NumberOfAccounts returns the number of accounts in the blockchain.
func (mf *MetaFetcher) NumberOfAccounts() (uint64, error) {
	resp, err := http.Get(mf.url)
	if err != nil {
		mf.log.Errorf("failed to get number of accounts: %v", err)
		return 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			mf.log.Errorf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		mf.log.Errorf("failed to read response body: %v", err)
		return 0, err
	}

	numStr := strings.TrimSpace(string(body))
	num, err := strconv.Atoi(numStr)
	if err != nil {
		mf.log.Errorf("failed to convert number of accounts: %v", err)
		return 0, err
	}

	return uint64(num), nil
}

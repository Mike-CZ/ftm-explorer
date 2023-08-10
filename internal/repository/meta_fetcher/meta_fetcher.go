package meta_fetcher

import (
	"encoding/json"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// MetaFetcher represents a meta fetcher. It is responsible for fetching the blockchain metadata.
type MetaFetcher struct {
	log                   logger.ILogger
	numberOfAccountsUrl   string
	diskSizePer100MTxsUrl string
	timeToFinalityUrl     string
}

// NewMetaFetcher returns a new meta fetcher.
func NewMetaFetcher(cfg *config.MetaFetcher, log logger.ILogger) *MetaFetcher {
	return &MetaFetcher{
		log:                   log.ModuleLogger("meta_fetcher"),
		numberOfAccountsUrl:   cfg.NumberOfAccountsUrl,
		diskSizePer100MTxsUrl: cfg.DiskSizePer100MTxsUrl,
		timeToFinalityUrl:     cfg.TimeToFinalityUrl,
	}
}

// NumberOfAccounts returns the number of accounts in the blockchain.
func (mf *MetaFetcher) NumberOfAccounts() (uint64, error) {
	resp, err := http.Get(mf.numberOfAccountsUrl)
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

// DiskSizePer100MTxs returns the disk size per 100M transactions.
func (mf *MetaFetcher) DiskSizePer100MTxs() (uint64, error) {
	resp, err := http.Get(mf.diskSizePer100MTxsUrl)
	if err != nil {
		mf.log.Errorf("failed to get disk size per 100m txs: %v", err)
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
		mf.log.Errorf("failed to convert disk size per 100m txs: %v", err)
		return 0, err
	}

	return uint64(num), nil
}

// TimeToFinality returns the time to finality in the blockchain.
func (mf *MetaFetcher) TimeToFinality() (float64, error) {
	resp, err := http.Get(mf.timeToFinalityUrl)
	if err != nil {
		mf.log.Errorf("failed to get time to finality: %v", err)
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

	res := struct {
		TimeToFinality float64 `json:"ttf"`
	}{}

	if err := json.Unmarshal(body, &res); err != nil {
		mf.log.Errorf("failed to unmarshal response body: %v", err)
		return 0, err
	}

	return res.TimeToFinality, nil
}

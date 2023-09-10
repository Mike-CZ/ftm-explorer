package faucet

import (
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"testing"

	"github.com/golang/mock/gomock"
)

const (
	kClaimLimitSeconds = 60
	kClaimTokensAmount = 1.5
)

// Test that the new tokens request can be created.
func TestFaucet_RequestTokens(t *testing.T) {
	faucet, repo := createFaucet(t)
	ipAddress := "192.168.0.1"

	// expect a call to the repository to get the latest tokens request
	repo.EXPECT().GetLatestTokensRequest(ipAddress).Return(nil, nil)

	// expect a call to the repository to add a new tokens request
	repo.EXPECT().AddTokensRequest(gomock.Any()).Return(nil)

	tr, err := faucet.RequestTokens(ipAddress)
	if err != nil {
		t.Fatalf("RequestTokens failed: %v", err)
	}

	// validate the returned tokens request
	if tr == nil {
		t.Fatal("RequestTokens returned nil")
	}
	if tr.IpAddress != ipAddress {
		t.Errorf("RequestTokens returned wrong ip address: %s", tr.IpAddress)
	}
	if len(tr.Phrase) == 0 {
		t.Errorf("RequestTokens returned empty phrase")
	}
	if tr.Receiver != nil {
		t.Errorf("RequestTokens returned wrong receiver")
	}
	if tr.ClaimedAt != nil {
		t.Errorf("RequestTokens returned wrong claim time")
	}
}

// createFaucet creates a new faucet instance for testing.
func createFaucet(t *testing.T) (*Faucet, *repository.MockRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepository(ctrl)
	cfg := &config.Faucet{
		ClaimLimitSeconds: kClaimLimitSeconds,
		ClaimTokensAmount: kClaimTokensAmount,
	}
	return NewFaucet(mockRepository, cfg, logger.NewMockLogger()), mockRepository
}

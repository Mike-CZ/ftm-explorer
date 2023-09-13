package faucet

import (
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

const (
	kClaimLimitSeconds = 60
	kClaimTokensAmount = 1.5
)

// Test that the new tokens request can be created.
func TestFaucet_RequestTokens(t *testing.T) {
	faucet, pg, repo := createFaucet(t)
	ipAddress := "192.168.0.1"

	// expect a call to the repository to get the latest tokens request
	repo.EXPECT().GetLatestTokensRequest(ipAddress).Return(nil, nil)

	// expect a call to the phrase generator to generate a new phrase
	pg.EXPECT().GeneratePhrase().Return("test-phrase")

	// expect a call to the repository to add a new tokens request
	repo.EXPECT().AddTokensRequest(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    "test-phrase",
	}).Return(nil)

	phrase, err := faucet.RequestTokens(ipAddress)
	if err != nil {
		t.Fatalf("RequestTokens failed: %v", err)
	}

	if phrase != challengePrefix+"test-phrase" {
		t.Fatalf("Invalid phrase returned: %s", phrase)
	}
}

// Test that the existing tokens request is returned.
func TestFaucet_RequestTokensAlreadyPending(t *testing.T) {
	faucet, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"

	tr := &types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    "pending-phrase",
	}

	// expect a call to the repository to get the latest tokens request
	repo.EXPECT().GetLatestTokensRequest(ipAddress).Return(tr, nil)

	// the request defined above should be returned, because it is pending
	_, err := faucet.RequestTokens(ipAddress)
	if err != nil {
		t.Fatalf("RequestTokens failed: %v", err)
	}

}

// Test that error is returned when claim limit is not reached
func TestFaucet_RequestTokensClaimLimitNotReached(t *testing.T) {
	faucet, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	receiver := common.Address{0x01}

	now := time.Now().Unix()
	tr := &types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    "pending-phrase",
		Receiver:  &receiver,
		ClaimedAt: &now,
	}

	// expect a call to the repository to get the latest tokens request
	repo.EXPECT().GetLatestTokensRequest(ipAddress).Return(tr, nil)

	// we should get an error, because the claim limit is not reached
	_, err := faucet.RequestTokens(ipAddress)
	if err == nil || !strings.Contains(err.Error(), "must wait") {
		t.Fatal("RequestTokens did not return error")
	}
}

// Test that request can be made when claim limit is reached
func TestFaucet_RequestTokensClaimLimitReached(t *testing.T) {
	faucet, pg, repo := createFaucet(t)
	ipAddress := "192.168.0.1"

	lastClaim := time.Now().Add(-time.Duration(kClaimLimitSeconds) * time.Second).Unix()
	tr := &types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    "pending-phrase",
		ClaimedAt: &lastClaim,
	}

	// expect a call to the repository to get the latest tokens request
	repo.EXPECT().GetLatestTokensRequest(ipAddress).Return(tr, nil)

	// expect a call to the phrase generator to generate a new phrase
	pg.EXPECT().GeneratePhrase().Return("different-phrase")

	// expect a call to the repository to add a new tokens request
	repo.EXPECT().AddTokensRequest(gomock.Any()).Return(nil)

	// we should get a new tokens request, because the claim limit is reached
	phrase, err := faucet.RequestTokens(ipAddress)
	if err != nil {
		t.Fatalf("RequestTokens failed: %v", err)
	}

	if phrase != challengePrefix+"different-phrase" {
		t.Fatalf("Invalid phrase returned: %s", phrase)
	}
}

// Test that the tokens can be claimed.
func TestFaucet_ClaimTokens(t *testing.T) {
	faucet, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestTokensRequest(ipAddress).Return(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
	}, nil)

	// expect a call to the repository to update the tokens request
	now := time.Now().Unix()
	repo.EXPECT().UpdateTokensRequest(gomock.Eq(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
		Receiver:  &receiver,
		ClaimedAt: &now,
	})).Return(nil)

	// TODO: Expect call to tokens transfer

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, challengePrefix+phrase, receiver)
	if err != nil {
		t.Fatalf("ClaimTokens failed: %v", err)
	}
}

// Test that error is returned when tokens request is not found.
func TestFaucet_ClaimTokensNoPendingRequest(t *testing.T) {
	faucet, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestTokensRequest(ipAddress).Return(nil, nil)

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, challengePrefix+phrase, receiver)
	if err == nil || !strings.Contains(err.Error(), "no request found") {
		t.Fatal("ClaimTokens did not return error")
	}
}

// Test that error is returned when phrase does not match.
func TestFaucet_ClaimTokensPhraseMismatch(t *testing.T) {
	faucet, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestTokensRequest(ipAddress).Return(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
	}, nil)

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, challengePrefix+"different-phrase", receiver)
	if err == nil || err.Error() != "invalid phrase" {
		t.Fatal("ClaimTokens did not return error")
	}
}

// Test that error is returned when tokens are already claimed.
func TestFaucet_ClaimTokensAlreadyClaimed(t *testing.T) {
	faucet, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	claimed := time.Now().Unix()
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestTokensRequest(ipAddress).Return(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
		Receiver:  &receiver,
		ClaimedAt: &claimed,
	}, nil)

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, challengePrefix+phrase, receiver)
	if err == nil || err.Error() != "tokens already claimed" {
		t.Fatal("ClaimTokens did not return error")
	}
}

// Test that error is returned when prefix is not present.
func TestFaucet_ClaimTokensNoPrefix(t *testing.T) {
	faucet, _, _ := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, phrase, receiver)
	if err == nil || err.Error() != "invalid phrase" {
		t.Fatal("ClaimTokens did not return error")
	}
}

// createFaucet creates a new faucet instance for testing.
func createFaucet(t *testing.T) (*Faucet, *MockFaucetPhraseGenerator, *repository.MockRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepository(ctrl)
	mockPhraseGenerator := NewMockFaucetPhraseGenerator(ctrl)
	cfg := &config.Faucet{
		ClaimLimitSeconds: kClaimLimitSeconds,
		ClaimTokensAmount: kClaimTokensAmount,
	}
	return NewFaucet(mockPhraseGenerator, mockRepository, cfg), mockPhraseGenerator, mockRepository
}

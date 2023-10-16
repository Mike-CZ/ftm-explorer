package faucet

import (
	"fmt"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/types"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

const (
	kClaimLimitSeconds = 60
	kClaimTokensAmount = 1.5
	kErc20Address      = "0x3bc666c4073853a59a7bfb0184298551d922f1df"
	kErc20MintAmount   = 10_000
)

// Test that the new tokens request can be created.
func TestFaucet_RequestTokens(t *testing.T) {
	faucet, pg, _, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"

	// expect a call to the repository to get the latest tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(nil, nil)

	// expect a call to the repository to get the latest claimed tokens requests
	repo.EXPECT().GetLatestClaimedTokensRequests(ipAddress, gomock.Any()).Return([]types.TokensRequest{}, nil)

	// expect a call to the phrase generator to generate a new phrase
	pg.EXPECT().GeneratePhrase().Return("test-phrase", nil)

	// expect a call to the repository to add a new tokens request
	repo.EXPECT().AddTokensRequest(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    "test-phrase",
	}).Return(nil)

	phrase, err := faucet.RequestTokens(ipAddress)
	if err != nil {
		t.Fatalf("RequestTokens failed: %v", err)
	}

	if phrase != kFaucetChallengePrefix+"test-phrase" {
		t.Fatalf("Invalid phrase returned: %s", phrase)
	}
}

// Test that the existing tokens request is returned.
func TestFaucet_RequestTokensAlreadyPending(t *testing.T) {
	faucet, _, _, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"

	tr := &types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    "pending-phrase",
	}

	// expect a call to the repository to get the latest tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(tr, nil)

	// the request defined above should be returned, because it is pending
	_, err := faucet.RequestTokens(ipAddress)
	if err != nil {
		t.Fatalf("RequestTokens failed: %v", err)
	}

}

// Test that error is returned when claim limit is not reached
func TestFaucet_RequestTokensClaimLimitNotReached(t *testing.T) {
	faucet, _, _, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	receiver := common.Address{0x01}

	now := time.Now().Unix() - 12*60*60 // 12 hours ago
	trs := []types.TokensRequest{
		{
			IpAddress: ipAddress,
			Phrase:    "pending-phrase 1",
			Receiver:  &receiver,
			ClaimedAt: &now,
		},
		{
			IpAddress: ipAddress,
			Phrase:    "pending-phrase 2",
			Receiver:  &receiver,
			ClaimedAt: &now,
		},
		{
			IpAddress: ipAddress,
			Phrase:    "pending-phrase 3",
			Receiver:  &receiver,
			ClaimedAt: &now,
		},
	}

	// expect a call to the repository to get the latest tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(nil, nil)

	// expect a call to the repository to get the latest claimed tokens requests
	repo.EXPECT().GetLatestClaimedTokensRequests(ipAddress, gomock.Any()).Return(trs, nil)

	// we should get an error, because the claim limit is not reached
	_, err := faucet.RequestTokens(ipAddress)
	if err == nil || !strings.Contains(err.Error(), "too many requests") {
		t.Fatal("RequestTokens did not return error")
	}
}

// Test that request can be made when claim limit is reached
func TestFaucet_RequestTokensClaimLimitReached(t *testing.T) {
	faucet, pg, _, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"

	// expect a call to the repository to get the latest tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(nil, nil)

	// expect a call to the repository to get the latest claimed tokens requests
	repo.EXPECT().GetLatestClaimedTokensRequests(ipAddress, gomock.Any()).Return([]types.TokensRequest{}, nil)

	// expect a call to the phrase generator to generate a new phrase
	pg.EXPECT().GeneratePhrase().Return("different-phrase", nil)

	// expect a call to the repository to add a new tokens request
	repo.EXPECT().AddTokensRequest(gomock.Any()).Return(nil)

	// we should get a new tokens request, because the claim limit is reached
	phrase, err := faucet.RequestTokens(ipAddress)
	if err != nil {
		t.Fatalf("RequestTokens failed: %v", err)
	}

	if phrase != kFaucetChallengePrefix+"different-phrase" {
		t.Fatalf("Invalid phrase returned: %s", phrase)
	}
}

// Test that the tokens can be claimed.
func TestFaucet_ClaimTokens(t *testing.T) {
	faucet, _, wallet, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(&types.TokensRequest{
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

	// expect a call to wallet to send tokens
	wallet.EXPECT().SendWeiToAddress(gomock.Eq(getTokensAmountInWei(kClaimTokensAmount)), gomock.Eq(receiver)).Return(nil)

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, kFaucetChallengePrefix+phrase, receiver, nil)
	if err != nil {
		t.Fatalf("ClaimTokens failed: %v", err)
	}
}

// Test that error is returned when tokens request is not found.
func TestFaucet_ClaimTokensNoPendingRequest(t *testing.T) {
	faucet, _, _, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(nil, nil)

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, kFaucetChallengePrefix+phrase, receiver, nil)
	if err == nil || !strings.Contains(err.Error(), "no request found") {
		t.Fatal("ClaimTokens did not return error")
	}
}

// Test that error is returned when phrase does not match.
func TestFaucet_ClaimTokensPhraseMismatch(t *testing.T) {
	faucet, _, _, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
	}, nil)

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, kFaucetChallengePrefix+"different-phrase", receiver, nil)
	if err == nil || err.Error() != "invalid phrase" {
		t.Fatal("ClaimTokens did not return error")
	}
}

// Test that error is returned when tokens are already claimed.
func TestFaucet_ClaimTokensAlreadyClaimed(t *testing.T) {
	faucet, _, _, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	claimed := time.Now().Unix()
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
		Receiver:  &receiver,
		ClaimedAt: &claimed,
	}, nil)

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, kFaucetChallengePrefix+phrase, receiver, nil)
	if err == nil || err.Error() != "tokens already claimed" {
		t.Fatal("ClaimTokens did not return error")
	}
}

// Test that error is returned when prefix is not present.
func TestFaucet_ClaimTokensNoPrefix(t *testing.T) {
	faucet, _, _, _, _ := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, phrase, receiver, nil)
	if err == nil || err.Error() != "invalid phrase" {
		t.Fatal("ClaimTokens did not return error")
	}
}

// Test that error is returned when wallet returns error and claim is reset.
func TestFaucet_ClaimTokensWalletError(t *testing.T) {
	faucet, _, wallet, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
	}, nil)

	repo.EXPECT().UpdateTokensRequest(gomock.Any()).Return(nil)

	// expect a call to wallet to send tokens
	wallet.EXPECT().SendWeiToAddress(gomock.Eq(getTokensAmountInWei(kClaimTokensAmount)), gomock.Eq(receiver)).Return(fmt.Errorf("error sending wei"))

	// expect call to reset the claim
	repo.EXPECT().UpdateTokensRequest(gomock.Eq(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
		ClaimedAt: nil,
		Receiver:  nil,
	})).Return(nil)

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, kFaucetChallengePrefix+phrase, receiver, nil)
	if err == nil {
		t.Fatalf("ClaimTokens did not return error")
	}
}

// Test that error is returned when erc20 is unknown.
func TestFaucet_ClaimErc20TokensUnknownAddress(t *testing.T) {
	faucet, _, _, _, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(&types.TokensRequest{
		IpAddress: ipAddress,
		Phrase:    phrase,
	}, nil)

	// claim tokens with unknown address
	addr := common.Address{0x02}
	err := faucet.ClaimTokens(ipAddress, kFaucetChallengePrefix+phrase, receiver, &addr)
	if err == nil || err.Error() != "unknown erc20 token" {
		t.Fatal("ClaimTokens did not return error")
	}
}

// Test that the erc 20 tokens can be minted.
func TestFaucet_MintErc20Tokens(t *testing.T) {
	faucet, _, _, erc20Wallet, repo := createFaucet(t)
	ipAddress := "192.168.0.1"
	phrase := "test-phrase"
	receiver := common.Address{0x01}

	// expect a call to the repository to get the tokens request
	repo.EXPECT().GetLatestUnclaimedTokensRequest(ipAddress).Return(&types.TokensRequest{
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

	// expect a call to wallet to mint tokens
	addr := common.HexToAddress(kErc20Address)
	erc20Wallet.EXPECT().MintErc20TokensToAddress(gomock.Eq(addr), gomock.Eq(receiver), gomock.Eq(new(big.Int).SetUint64(kErc20MintAmount))).Return(nil)

	// claim tokens
	err := faucet.ClaimTokens(ipAddress, kFaucetChallengePrefix+phrase, receiver, &addr)
	if err != nil {
		t.Fatalf("ClaimTokens failed: %v", err)
	}
}

// test that the amount of tokens is converted to wei correctly.
func TestFaucet_GetTokensAmountInWei(t *testing.T) {
	wei := getTokensAmountInWei(0.5)
	if wei.Cmp(big.NewInt(500_000_000_000_000_000)) != 0 {
		t.Fatal("Invalid wei amount")
	}

	wei = getTokensAmountInWei(.000_000_000_000_000_001)
	if wei.Cmp(big.NewInt(1)) != 0 {
		t.Fatal("Invalid wei amount")
	}
}

// createFaucet creates a new faucet instance for testing.
func createFaucet(t *testing.T) (*Faucet, *MockFaucetPhraseGenerator, *MockFaucetWallet, *MockFaucetWallet, *repository.MockRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepository := repository.NewMockRepository(ctrl)
	mockPhraseGenerator := NewMockFaucetPhraseGenerator(ctrl)
	mockWallet := NewMockFaucetWallet(ctrl)
	cfg := &config.Faucet{
		ClaimLimitSeconds: kClaimLimitSeconds,
		ClaimTokensAmount: kClaimTokensAmount,
	}
	mockErc20Wallet := NewMockFaucetWallet(ctrl)
	erc20s := []FaucetErc20{
		{
			address:    common.HexToAddress(kErc20Address),
			wallet:     mockErc20Wallet,
			mintAmount: new(big.Int).SetUint64(kErc20MintAmount),
		},
	}
	return NewFaucet(cfg, mockPhraseGenerator, mockWallet, erc20s, mockRepository), mockPhraseGenerator, mockWallet, mockErc20Wallet, mockRepository
}

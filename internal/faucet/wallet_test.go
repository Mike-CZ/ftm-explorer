package faucet

import (
	"context"
	"fmt"
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"ftm-explorer/internal/repository/db"
	"ftm-explorer/internal/repository/rpc"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/mock/gomock"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Test that the faucet wallet can send wei to an address.
func TestFaucetWallet_SendWeiToAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	receiver := common.HexToAddress("0x5A4b203939F9757A703e009fA9B733Cf33d5821b")

	// start test chain
	url := startTestChain(t)

	// initialize stubs and rpc client that will connect to test chain
	log := logger.NewMockLogger()
	database := db.NewMockDatabase(ctrl)
	client, err := rpc.NewOperaRpc(&config.Rpc{
		OperaRpcUrl: url,
	})
	if err != nil {
		t.Fatal(err)
	}

	// initialize eth client so we can check balances of accounts
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	// initialize repository
	repo := repository.NewRepository(10_000, client, database, nil)

	// initialize wallet
	wallet, err := NewWallet(repo, log, "bb39aa88008bc6260ff9ebc816178c47a01c44efe55810ea1f271c00f5878812")
	if err != nil {
		t.Fatal(err)
	}

	// get initial balances
	senderBalance, err := ethClient.BalanceAt(context.Background(), wallet.from, nil)
	if err != nil {
		t.Fatal(err)
	}
	receiverBalance, err := ethClient.BalanceAt(context.Background(), receiver, nil)
	if err != nil {
		t.Fatal(err)
	}

	// send 2.5 eth to receiver
	wei := getTokensAmountInWei(2.5)
	err = wallet.SendWeiToAddress(wei, receiver)

	// get updated balances
	senderBalanceUpdated, err := ethClient.BalanceAt(context.Background(), wallet.from, nil)
	if err != nil {
		t.Fatal(err)
	}
	receiverBalanceUpdated, err := ethClient.BalanceAt(context.Background(), receiver, nil)
	if err != nil {
		t.Fatal(err)
	}

	// calculate tx fee, 21_000 is simple tx cost, 6_721_975 is gas price configured in ganache
	txFee := new(big.Int).Mul(new(big.Int).SetUint64(21_000), new(big.Int).SetUint64(6_721_975))

	// check amount was transferred
	if senderBalanceUpdated.Cmp(senderBalance.Sub(senderBalance, wei).Sub(senderBalance, txFee)) != 0 {
		t.Fatalf("sender balance was not updated correctly")
	}
	if receiverBalanceUpdated.Cmp(receiverBalance.Add(receiverBalance, wei)) != 0 {
		t.Fatalf("receiver balance was not updated correctly")
	}
}

// startTestChain starts a new test chain. It runs ganache inside docker container.
// It returns the url of the chain.
func startTestChain(t *testing.T) string {
	t.Helper()
	// setup db container
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	container, port := createContainer(t, ctx)
	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
	return fmt.Sprintf("http://localhost:%s", port)
}

// createContainer creates a test container for postgres database
func createContainer(t *testing.T, ctx context.Context) (testcontainers.Container, string) {
	t.Helper()
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "trufflesuite/ganache-cli:v6.12.2",
			ExposedPorts: []string{"8545/tcp"},
			Env:          nil,
			WaitingFor:   wait.ForListeningPort("8545/tcp"),
			Cmd: []string{
				"--chainId", strconv.Itoa(1337),
				"--gasLimit", strconv.Itoa(10_000_000_000),
				"--gasPrice", strconv.Itoa(6_721_975),
				// address: 0x9Cc2F0FD184E93049A9a6C6C63bc258A39D4B54D
				"--account", "0xbb39aa88008bc6260ff9ebc816178c47a01c44efe55810ea1f271c00f5878812,200000000000000000000",
				// address: 0x5A4b203939F9757A703e009fA9B733Cf33d5821b
				"--account", "0x29c8b4ff78e41dafd561f5cd4a90103faf20a5b509a4b6281947b8fcdcfa8f71,100000000000000000000",
			},
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		t.Fatalf("failed to create container: %v", err)
	}
	p, err := container.MappedPort(ctx, "8545/tcp")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}
	// wait for the chain to be ready
	time.Sleep(time.Second)

	return container, p.Port()
}

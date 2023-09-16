package faucet

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Test that the faucet wallet can send wei to an address.
func TestFaucetWallet_SendWeiToAddress(t *testing.T) {
	url := startTestChain(t)
	fmt.Println(url)
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
				"--account", "0xbb39aa88008bc6260ff9ebc816178c47a01c44efe55810ea1f271c00f5878812,100000000000000000000",
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

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
	"strings"
	"testing"
	"time"

	abi2 "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

func TestFaucetWallet_MintErc20Tokens(t *testing.T) {
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

	// initialize repository
	repo := repository.NewRepository(10_000, client, database, nil)

	// initialize wallet
	wallet, err := NewWallet(repo, log, "bb39aa88008bc6260ff9ebc816178c47a01c44efe55810ea1f271c00f5878812")
	if err != nil {
		t.Fatal(err)
	}

	// deploy erc20 contract
	contractAddress := deployErc20Contract(t, client)

	// mint erc 20 tokens
	err = wallet.MintErc20TokensToAddress(contractAddress, receiver, big.NewInt(1000))
	if err != nil {
		t.Fatal(err)
	}
}

// deployErc20Contract deploys an erc20 contract to the test chain.
// It returns the address of the contract.
func deployErc20Contract(t *testing.T, rpc rpc.IRpc) common.Address {
	t.Helper()

	pk, err := crypto.HexToECDSA("bb39aa88008bc6260ff9ebc816178c47a01c44efe55810ea1f271c00f5878812")
	if err != nil {
		t.Fatalf("failed to parse private key: %v", err)
	}

	nonce, err := rpc.PendingNonceAt(context.Background(), common.HexToAddress("0x9Cc2F0FD184E93049A9a6C6C63bc258A39D4B54D"))
	if err != nil {
		t.Fatalf("failed to get nonce: %v", err)
	}

	// definition of erc20 constructor abi
	definition := `[{"inputs": [{"internalType": "string","name": "_name","type": "string"},{"internalType": "string","name": "_symbol","type": "string"},{"internalType": "address","name": "minter","type": "address"}],"stateMutability": "nonpayable","type": "constructor"}]`
	abi, err := abi2.JSON(strings.NewReader(definition))
	if err != nil {
		t.Fatalf("failed to parse abi: %v", err)
	}
	// pack constructor params
	params, err := abi.Pack("", "erc20-test", "TST", common.HexToAddress("0x9Cc2F0FD184E93049A9a6C6C63bc258A39D4B54D"))
	if err != nil {
		t.Fatalf("failed to pack constructor params: %v", err)
	}
	// bytecode of erc20 contract
	bytecode := common.FromHex("0x60806040526005805460ff191660121790553480156200001e57600080fd5b5060405162000b8a38038062000b8a833981016040819052620000419162000166565b60058054610100600160a81b0319163361010002179055600362000066848262000282565b50600462000075838262000282565b506001600160a01b03166000908152600660205260409020805460ff19166001179055506200034e9050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620000c957600080fd5b81516001600160401b0380821115620000e657620000e6620000a1565b604051601f8301601f19908116603f01168101908282118183101715620001115762000111620000a1565b816040528381526020925086838588010111156200012e57600080fd5b600091505b8382101562000152578582018301518183018401529082019062000133565b600093810190920192909252949350505050565b6000806000606084860312156200017c57600080fd5b83516001600160401b03808211156200019457600080fd5b620001a287838801620000b7565b94506020860151915080821115620001b957600080fd5b50620001c886828701620000b7565b604086015190935090506001600160a01b0381168114620001e857600080fd5b809150509250925092565b600181811c908216806200020857607f821691505b6020821081036200022957634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200027d57600081815260208120601f850160051c81016020861015620002585750805b601f850160051c820191505b81811015620002795782815560010162000264565b5050505b505050565b81516001600160401b038111156200029e576200029e620000a1565b620002b681620002af8454620001f3565b846200022f565b602080601f831160018114620002ee5760008415620002d55750858301515b600019600386901b1c1916600185901b17855562000279565b600085815260208120601f198616915b828110156200031f57888601518255948401946001909101908401620002fe565b50858210156200033e5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61082c806200035e6000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806370a082311161008c5780639b19251a116100665780639b19251a146101cb578063a9059cbb146101ee578063d77ecb4114610201578063dd62ed3e1461021457600080fd5b806370a08231146101735780638da5cb5b1461019357806395d89b41146101c357600080fd5b806306fdde03146100d4578063095ea7b3146100f257806318160ddd1461011557806323b872dd1461012c578063313ce5671461013f57806340c10f191461015e575b600080fd5b6100dc61023f565b6040516100e9919061065b565b60405180910390f35b6101056101003660046106c5565b6102cd565b60405190151581526020016100e9565b61011e60005481565b6040519081526020016100e9565b61010561013a3660046106ef565b61033a565b60055461014c9060ff1681565b60405160ff90911681526020016100e9565b61017161016c3660046106c5565b61043d565b005b61011e61018136600461072b565b60016020526000908152604090205481565b6005546101ab9061010090046001600160a01b031681565b6040516001600160a01b0390911681526020016100e9565b6100dc610533565b6101056101d936600461072b565b60066020526000908152604090205460ff1681565b6101056101fc3660046106c5565b610540565b61017161020f36600461072b565b6105d0565b61011e61022236600461074d565b600260209081526000928352604080842090915290825290205481565b6003805461024c90610780565b80601f016020809104026020016040519081016040528092919081815260200182805461027890610780565b80156102c55780601f1061029a576101008083540402835291602001916102c5565b820191906000526020600020905b8154815290600101906020018083116102a857829003601f168201915b505050505081565b3360008181526002602090815260408083206001600160a01b038716808552925280832085905551919290917f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925906103289086815260200190565b60405180910390a35060015b92915050565b3360009081526006602052604081205460ff1661038a576001600160a01b0384166000908152600260209081526040808320338452909152812080548492906103849084906107d0565b90915550505b6001600160a01b038416600090815260016020526040812080548492906103b29084906107d0565b90915550506001600160a01b038316600090815260016020526040812080548492906103df9084906107e3565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161042b91815260200190565b60405180910390a35060019392505050565b60055461010090046001600160a01b031633148061046a57503360009081526006602052604090205460ff165b6104a95760405162461bcd60e51b815260206004820152600b60248201526a1b9bdd08185b1b1bddd95960aa1b60448201526064015b60405180910390fd5b6001600160a01b038216600090815260016020526040812080548392906104d19084906107e3565b92505081905550806000808282546104e991906107e3565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b6004805461024c90610780565b336000908152600160205260408120805483919083906105619084906107d0565b90915550506001600160a01b0383166000908152600160205260408120805484929061058e9084906107e3565b90915550506040518281526001600160a01b0384169033907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef90602001610328565b60055461010090046001600160a01b03163314806105fd57503360009081526006602052604090205460ff165b6106375760405162461bcd60e51b815260206004820152600b60248201526a1b9bdd08185b1b1bddd95960aa1b60448201526064016104a0565b6001600160a01b03166000908152600660205260409020805460ff19166001179055565b600060208083528351808285015260005b818110156106885785810183015185820160400152820161066c565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b03811681146106c057600080fd5b919050565b600080604083850312156106d857600080fd5b6106e1836106a9565b946020939093013593505050565b60008060006060848603121561070457600080fd5b61070d846106a9565b925061071b602085016106a9565b9150604084013590509250925092565b60006020828403121561073d57600080fd5b610746826106a9565b9392505050565b6000806040838503121561076057600080fd5b610769836106a9565b9150610777602084016106a9565b90509250929050565b600181811c9082168061079457607f821691505b6020821081036107b457634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b81810381811115610334576103346107ba565b80820180821115610334576103346107ba56fea2646970667358221220d6ee47d24af850d1f79f5ce261911784725aa664bdafd0ebbab4bf0b8bf702c264736f6c63430008130033")

	// deploy contract
	tx := types.NewContractCreation(
		nonce, new(big.Int).SetUint64(0), 10_000_000_000, new(big.Int).SetUint64(6_721_975), append(bytecode, params...),
	)

	// sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1337)), pk)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}

	// send transaction
	if err = rpc.SendSignedTransaction(context.Background(), signedTx); err != nil {
		t.Fatalf("failed to send transaction: %v", err)
	}

	// get full transaction
	nTx, err := rpc.TransactionByHash(context.Background(), signedTx.Hash())
	if err != nil {
		t.Fatalf("failed to get transaction: %v", err)
	}

	// check if contract address is not nil
	if nTx.ContractAddress == nil {
		t.Fatalf("contract address is nil")
	}

	return *nTx.ContractAddress
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

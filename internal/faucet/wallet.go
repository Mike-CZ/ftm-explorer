package faucet

import (
	"crypto/ecdsa"
	"fmt"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"math/big"
	"strings"

	abi2 "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Wallet represents a faucet wallet. It is used to send wei to the given address.
type Wallet struct {
	repo repository.IRepository
	log  logger.ILogger
	pk   *ecdsa.PrivateKey
	from common.Address
}

// NewWallet returns a new wallet.
func NewWallet(repo repository.IRepository, log logger.ILogger, pk string) (*Wallet, error) {
	// initialize logger
	wl := log.ModuleLogger("faucet_wallet")

	// initialize private key
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	// get from address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Wallet{
		repo: repo,
		log:  wl,
		pk:   privateKey,
		from: fromAddress,
	}, nil
}

// SendWeiToAddress sends wei to the given address.
func (w *Wallet) SendWeiToAddress(amount *big.Int, receiver common.Address) error {
	// get nonce
	nonce, err := w.repo.PendingNonceAt(w.from)
	if err != nil {
		w.log.Criticalf("error getting nonce: %v", err)
	}

	// get gas price
	gasPrice, err := w.repo.SuggestGasPrice()
	if err != nil {
		w.log.Criticalf("error getting gas price: %v", err)
		return err
	}

	// get network id
	chainID, err := w.repo.NetworkID()
	if err != nil {
		w.log.Criticalf("error getting network id: %v", err)
		return err
	}

	// create transaction, set gas limit to 21000, which is the cost of a normal transaction
	tx := types.NewTransaction(nonce, receiver, amount, 21_000, gasPrice, nil)

	// sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), w.pk)
	if err != nil {
		w.log.Criticalf("error signing transaction: %v", err)
		return err
	}

	// send transaction
	if err = w.repo.SendSignedTransaction(signedTx); err != nil {
		w.log.Criticalf("error sending transaction: %v", err)
		return err
	}

	return nil
}

// MintErc20TokensToAddress mints erc20 tokens to the given address.
func (w *Wallet) MintErc20TokensToAddress(contract common.Address, receiver common.Address, amount *big.Int) error {
	// get nonce
	nonce, err := w.repo.PendingNonceAt(w.from)
	if err != nil {
		w.log.Criticalf("error getting nonce: %v", err)
	}

	// get gas price
	gasPrice, err := w.repo.SuggestGasPrice()
	if err != nil {
		w.log.Criticalf("error getting gas price: %v", err)
		return err
	}

	// get network id
	chainID, err := w.repo.NetworkID()
	if err != nil {
		w.log.Criticalf("error getting network id: %v", err)
		return err
	}

	// get erc20 mint data
	data, err := getErc20MintData(receiver, amount)
	if err != nil {
		w.log.Criticalf("error getting erc20 mint data: %v", err)
		return err
	}

	// create transaction
	tx := types.NewTransaction(nonce, contract, new(big.Int).SetUint64(0), 100_000, gasPrice, data)

	// sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), w.pk)
	if err != nil {
		w.log.Criticalf("error signing transaction: %v", err)
		return err
	}

	// send transaction
	if err = w.repo.SendSignedTransaction(signedTx); err != nil {
		w.log.Criticalf("error sending transaction: %v", err)
		return err
	}

	return nil
}

// getErc20MintData returns the erc20 mint data.
func getErc20MintData(receiver common.Address, amount *big.Int) ([]byte, error) {
	definition := `[{"inputs":[{"internalType": "address","name": "recipient","type": "address"},{"internalType": "uint256","name": "amount","type": "uint256"}],"name": "mint","outputs": [],"stateMutability": "nonpayable","type": "function"}]`
	abi, err := abi2.JSON(strings.NewReader(definition))
	if err != nil {
		return nil, err
	}
	// pack constructor params
	data, err := abi.Pack("mint", receiver, amount)
	if err != nil {
		return nil, err
	}
	return data, nil
}

package faucet

import (
	"crypto/ecdsa"
	"fmt"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"math/big"

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

	// build transaction data
	// identifier is first 4 bytes of keccak256 hash of "mint(address,uint256)"
	var r [32]byte
	copy(r[32-len(receiver.Bytes()):], receiver.Bytes())
	var a [32]byte
	copy(a[32-len(amount.Bytes()):], amount.Bytes())
	data := append([]byte{0x40, 0xc1, 0x0f, 0x19}, append(r[:], a[:]...)...)
	
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

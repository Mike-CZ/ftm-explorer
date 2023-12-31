package auth

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// VerifySignature verifies signature generated using web3.eth.personal.sign Metamask call
// signature documented at:
// https://geth.ethereum.org/docs/rpc/ns-personal#personal_sign
func VerifySignature(message string, address common.Address, signature []byte) (bool, error) {
	recoveredAddr, err := ecRecover([]byte(message), signature)
	if err != nil {
		return false, err
	}
	return bytes.Equal(recoveredAddr.Bytes(), address.Bytes()), nil
}

// SignMessage signs the message using the provided private key
// Primarily used for testing purposes
func SignMessage(message string, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	signature, err := crypto.Sign(accounts.TextHash([]byte(message)), privateKey)
	if err != nil {
		return nil, err
	}
	signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper

	return signature, nil
}

// ecRecover returns the address for the account that was used to create the signature.
// copy of internal go-ethereum function:
// https://github.com/ethereum/go-ethereum/blob/v1.10.9/internal/ethapi/api.go#L524
func ecRecover(data []byte, sig []byte) (common.Address, error) {
	if len(sig) != crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.SigToPub(accounts.TextHash(data), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*rpk), nil
}

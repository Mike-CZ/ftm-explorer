package auth

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

// Test that signature verification works
func TestVerification_VerifySignature(t *testing.T) {
	// Decode the private key from hex
	privateKey, err := crypto.HexToECDSA("bb39aa88008bc6260ff9ebc816178c47a01c44efe55810ea1f271c00f5878812")
	if err != nil {
		t.Fatal(err)
	}

	// Derive the Ethereum address from the private key
	derivedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// sign the message
	message := "\"Sign following challenge:\n test"
	signature, err := SignMessage(message, privateKey)

	result, err := VerifySignature(message, derivedAddress, signature)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Fatal("signature verification failed")
	}
}

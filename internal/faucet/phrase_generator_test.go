package faucet

import (
	"testing"
)

// Test that the phrase generator returns a valid phrase.
func TestMockFaucetPhraseGenerator_GeneratePhrase(t *testing.T) {
	generator := NewPhraseGenerator()
	phrase, err := generator.GeneratePhrase()
	if err != nil {
		t.Fatalf("GeneratePhrase failed: %v", err)
	}
	if len(phrase) == 0 {
		t.Fatalf("Invalid phrase returned: %s", phrase)
	}
}

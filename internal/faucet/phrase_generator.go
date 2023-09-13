package faucet

// PhraseGenerator is an implementation of the IFaucetPhraseGenerator interface.
// It is used to generate a phrase for the faucet.
type PhraseGenerator struct{}

// NewPhraseGenerator creates a new PhraseGenerator.
func NewPhraseGenerator() *PhraseGenerator {
	return &PhraseGenerator{}
}

// GeneratePhrase generates a phrase for the faucet.
func (f *PhraseGenerator) GeneratePhrase() string {
	return "some-phrase"
}

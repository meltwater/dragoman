package cryptography

import "fmt"

type WildcardDecryptionStrategy struct {
	Strategies map[string]Decryptor
}

type StrategyBuilder func() (Decryptor, error)

// NewWildcardDecryptionStrategy is the initializer function for WildcardDecryptionStrategy
func NewWildcardDecryptionStrategy(builders []StrategyBuilder) (*WildcardDecryptionStrategy, error) {
	wcStrat := &WildcardDecryptionStrategy{
		Strategies: make(map[string]Decryptor),
	}

	for _, builder := range builders {
		strat, err := builder()
		if err != nil {
			return nil, err
		}

		wcStrat.Add(strat.Key(), strat)
	}

	return wcStrat, nil
}

func (wds WildcardDecryptionStrategy) Key() string {
	return "*"
}

// Decrypt will process the input string for the correct strategy and run decrypt on that strategy
func (wds WildcardDecryptionStrategy) Decrypt(input string) ([]byte, error) {
	// Figure out what the encryption strategy was
	etype := ExtractEncryptionType(input)

	var strategy Decryptor
	var exists bool

	if strategy, exists = wds.Strategies[etype]; !exists {
		return nil, fmt.Errorf("not configured for decrypting ENC[%s,...] values", etype)
	}

	return strategy.Decrypt(input)
}

// Add is a builder function to build up any applicable decryption strategies
func (wds *WildcardDecryptionStrategy) Add(key string, strategy Decryptor) *WildcardDecryptionStrategy {
	wds.Strategies[key] = strategy

	return wds
}

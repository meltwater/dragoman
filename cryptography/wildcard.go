package cryptography

import "fmt"

type WildcardDecryptionStrategy struct {
	Strategies map[string]Decryptor
}

// NewWildcardDecryptionStrategy is the initializer function for WildcardDecryptionStrategy
func NewWildcardDecryptionStrategy() *WildcardDecryptionStrategy {
	return &WildcardDecryptionStrategy{
		Strategies: make(map[string]Decryptor),
	}
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

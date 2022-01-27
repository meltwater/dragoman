package cryptography

// Encryptor defines the contract for encrypting data
type Encryptor interface {
	Encrypt([]byte) (string, error)
}

// Decryptor defines the contract for decrypting data
type Decryptor interface {
	Decrypt(string) ([]byte, error)
}

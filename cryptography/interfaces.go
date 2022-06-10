package cryptography

// Encryptor defines the contract for encrypting data
type Encryptor interface {
	Key() string
	Encrypt([]byte) (string, error)
}

// Decryptor defines the contract for decrypting data
type Decryptor interface {
	Key() string
	Decrypt(string) ([]byte, error)
}

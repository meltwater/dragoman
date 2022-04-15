package cmd

import (
	"fmt"

	"github.com/meltwater/dragoman/cryptography"
)

func processSMEncrypt(cfg *encryptConfig) error {
	var err error

	if cfg.AwsRegion == "" {
		return fmt.Errorf("an aws region must be provided for Secrets Manager encryption")
	}

	var strategy *cryptography.SecretsManagerCryptoStrategy
	if strategy, err = cryptography.NewSecretsManagerCryptoStrategy(cfg.AwsRegion); err != nil {
		return fmt.Errorf("unable to create secrets manager crypto strategy: %v", err)
	}

	var envelope string

	var secretKey []byte = nil
	if cfg.SecretKey != "" {
		secretKey = []byte(cfg.SecretKey)
	}

	if envelope, err = strategy.Encrypt(secretKey, cfg.Key); err != nil {
		return fmt.Errorf("error encountered attempting secrets manager encryption: %v", err)
	}

	cfg.Out.Write([]byte(envelope))
	cfg.Out.Write([]byte("\n"))

	return nil
}

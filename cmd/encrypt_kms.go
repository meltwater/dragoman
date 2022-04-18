package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/meltwater/dragoman/cryptography"
)

const maxLineLength = 64

func processKmsEncrypt(cfg *encryptConfig) error {
	var input []byte
	var err error

	if input, err = ioutil.ReadAll(cfg.In); err != nil {
		return fmt.Errorf("unable to read input: %v", err)
	}

	if cfg.AwsRegion == "" {
		return fmt.Errorf("an aws region must be provided for KMS encryption")
	}

	var strategy *cryptography.KmsCryptoStrategy
	if strategy, err = cryptography.NewKmsCryptoStrategy(cfg.AwsRegion); err != nil {
		return fmt.Errorf("unable to create kms crypto strategy: %v", err)
	}

	var envelope string
	if envelope, err = strategy.Encrypt(input, cfg.Key); err != nil {
		return fmt.Errorf("error encountered attempting KMS encryption: %v", err)
	}

	if cfg.WrapLines {
		for i := 0; i < len(envelope); i += maxLineLength {
			cfg.Out.Write([]byte(envelope[i:min(i+maxLineLength, len(envelope))]))
			cfg.Out.Write([]byte("\n"))
		}
	} else {
		cfg.Out.Write([]byte(envelope))
		cfg.Out.Write([]byte("\n"))
	}

	return nil
}

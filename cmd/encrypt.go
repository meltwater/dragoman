/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/meltwater/dragoman/cryptography"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt the provided string (via standard in)",
	Long: `Encrypt strings using the provided encryption strategy.

Examples: 

Encrypt with KMS
"My string to encrypt" | dragoman encrypt --kms-key-id myKmsKey`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Is this a kms encryption?
		var kmsKey string
		if kmsKey, err = cmd.Flags().GetString("kms-key-id"); err != nil {
			panic("No KMS Key provided for command encrypt. Please specify --kms-key-id")
		} else if kmsKey != "" {
			var (
				awsRegion string
				wrapLines bool
				err       error
			)
			// Get any other relevant flags or environment variables
			if awsRegion, err = cmd.Flags().GetString("aws-region"); err != nil {
				panic(err)
			}

			if wrapLines, err = cmd.Flags().GetBool("wrap"); err != nil {
				panic(err)
			}

			// Try and do the encryption
			if err = processKmsEncrypt(&kmsEncryptConfig{
				In:        os.Stdin,
				Out:       os.Stdout,
				KmsKey:    kmsKey,
				AwsRegion: awsRegion,
				WrapLines: wrapLines,
			}); err != nil {
				panic(err)
			}

			return
		}

		// Add other encryption methods here
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	// Setup Flags(this command only) and Persistent Flags (this command and sub commands)
	encryptCmd.Flags().String("kms-key-id", os.Getenv("KMS_KEY_ID"), "Provides the KMS Key ID")
	encryptCmd.Flags().String("aws-region", getFirstEnv("AWS_REGION", "AWS_DEFAULT_REGION"), "Provides the AWS region to use for KMS")
	encryptCmd.Flags().BoolP("wrap", "w", false, "Wrap long lines at 64 characters")
}

type kmsEncryptConfig struct {
	In        io.Reader
	Out       io.Writer
	KmsKey    string
	AwsRegion string
	WrapLines bool
}

const maxLineLength = 64

func processKmsEncrypt(cfg *kmsEncryptConfig) error {
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
	if envelope, err = strategy.Encrypt(input, cfg.KmsKey); err != nil {
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

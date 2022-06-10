package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt the provided string (via standard in)",
	Long: `Encrypt strings using the provided encryption strategy.

Examples: 

Encrypt with AWS KMS
"My string to encrypt" | dragoman encrypt --kms-key-id myKmsKey

Encrypt with AWS Secrets Manager
dragoman encrypt --sm-key-id mySecretsManagerKey --sm-secret-key myValuesKey`,
	Run: func(cmd *cobra.Command, args []string) {
		// KMS Envelope Encrpytion
		var kmsKey string
		if kmsKey, _ = cmd.Flags().GetString("kms-key-id"); kmsKey != "" {
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
			if err = processKmsEncrypt(&encryptConfig{
				In:        os.Stdin,
				Out:       os.Stdout,
				Key:       kmsKey,
				AwsRegion: awsRegion,
				WrapLines: wrapLines,
			}); err != nil {
				panic(err)
			}

			return
		}

		// Secrets Manager
		var smKey string
		if smKey, _ = cmd.Flags().GetString("sm-key-id"); smKey != "" {
			var (
				awsRegion string
				err       error
			)

			if awsRegion, err = cmd.Flags().GetString("aws-region"); err != nil {
				panic(err)
			}

			var smSecretKey, _ = cmd.Flags().GetString("sm-secret-key")
			if err = processSMEncrypt(&encryptConfig{
				Out:       os.Stdout,
				Key:       smKey,
				SecretKey: smSecretKey,
				AwsRegion: awsRegion,
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
	encryptCmd.Flags().String("sm-key-id", "", "Provides the Secrets Manager key to use")
	encryptCmd.Flags().String("sm-secret-key", "", "Provides the Key for Key/Value pairs in Secrets Manager")
	encryptCmd.Flags().String("aws-region", getFirstEnv("AWS_REGION", "AWS_DEFAULT_REGION"), "Provides the AWS region to use for KMS")
	encryptCmd.Flags().BoolP("wrap", "w", false, "Wrap long lines at 64 characters")
}

type encryptConfig struct {
	In        io.Reader
	Out       io.Writer
	Key       string
	SecretKey string // Secrets Manager specific
	AwsRegion string
	WrapLines bool
}

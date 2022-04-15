/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/meltwater/dragoman/cryptography"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt the provided string (via standard in or the file flag)",
	Long: `Automatically decrypt the string provided by standard in

The decryption strategy will be automatically detected`,
	Run: func(cmd *cobra.Command, args []string) {
		var input io.Reader = os.Stdin
		var output io.Writer = os.Stdout

		// File input
		if fname, err := cmd.Flags().GetString("input"); err == nil && fname != "" {
			var ferr error
			if input, ferr = os.Open(fname); ferr != nil {
				panic(fmt.Errorf("unable to open file \"%s\": %v", fname, ferr))
			}
		}

		// File output
		if fname, err := cmd.Flags().GetString("output"); err == nil && fname != "" {
			var ferr error
			if output, ferr = os.Create(fname); ferr != nil {
				panic(fmt.Errorf("unable to create output file \"%s\": %v", fname, ferr))
			}
		}

		// Be able to handle different encryption types
		strategy := cryptography.NewWildcardDecryptionStrategy()
		kmsStrategy, _ := cryptography.NewKmsCryptoStrategy("")
		smStrategy, _ := cryptography.NewSecretsManagerCryptoStrategy("")
		strategy.Add("KMS", kmsStrategy)
		strategy.Add("SECMAN", smStrategy)

		if err := processDecrypt(input, output, strategy); err != nil {
			panic(fmt.Errorf("unable to decrypt the provided text: %v", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringP("input", "i", "", "An optional input file to parse")
}

func processDecrypt(in io.Reader, out io.Writer, strategy cryptography.Decryptor) error {
	var (
		payload []byte
		err     error
		result  string
	)

	if payload, err = io.ReadAll(in); err != nil {
		return fmt.Errorf("unable to read input: %v", err)
	}

	if result, err = decryptEnvelopes(string(payload), strategy); err != nil {
		return fmt.Errorf("unable to decrypt input: %v", err)
	}

	out.Write([]byte(result))

	return nil
}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var VersionNumber string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dragoman",
	Short: "An encryption toolset",
	Long: `Dragoman is an encryption toolset that helps you encrypt
and decrypt your secrets. A common use case is when your secrets need
to live alongside your code.`,

	Run: func(cmd *cobra.Command, args []string) {
		vFlag, _ := cmd.Flags().GetBool("version")

		if vFlag {
			var output string
			if VersionNumber == "" {
				output = "Unknown Version"
			} else {
				output = VersionNumber
			}

			fmt.Println(output)

			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
}

package cmd

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/meltwater/dragoman/cryptography"
)

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func getFirstEnv(vars ...string) string {
	for _, envVar := range vars {
		rawVar := os.Getenv(envVar)
		if rawVar != "" {
			return rawVar
		}
	}

	return ""
}

// Strip whitespace from string
func stripWhitespace(a string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}

		return r
	}, a)
}

func decryptEnvelopes(input string, strategy cryptography.Decryptor) (output string, err error) {
	// Recover from any panics that happen at a lower level
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	replfn := func(envelope string) string {
		var data []byte
		var err error

		if data, err = strategy.Decrypt(stripWhitespace(envelope)); err != nil {
			panic(err)
		}

		return string(data)
	}

	output = cryptography.EnvelopeRegex.ReplaceAllStringFunc(input, replfn)

	return

}

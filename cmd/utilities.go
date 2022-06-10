package cmd

import (
	"os"
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

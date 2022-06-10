package cryptography

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	validStrategies = []string{
		"KMS",
		"SECMAN",
	}
	regexstr      = fmt.Sprintf("ENC\\[(%s),[a-zA-Z0-9+/=\\s]+\\]", strings.Join(validStrategies, "|"))
	EnvelopeRegex = regexp.MustCompile(regexstr)
)

// Converts a byte slice to a [32]byte as expected by NaCL
func AsNaCLKey(data []byte) (*[32]byte, error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("expected a 32 byte key for NaCL, got %d bytes", len(data))
	}

	var key [32]byte
	copy(key[:], data[0:32])
	return &key, nil
}

func WrapEncoding(key string, message []byte) string {
	return fmt.Sprintf("ENC[%s,%s]",
		key,
		base64.StdEncoding.EncodeToString(message))
}

func UnwrapEncoding(input string) ([]byte, error) {
	encoded := strings.Split(input[0:len(input)-1], ",")[1]
	message := make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))

	n, err := base64.StdEncoding.Decode(message, []byte(encoded))
	if err != nil {
		return nil, err
	}

	return message[0:n], nil
}

func ExtractEncryptionType(input string) string {
	submatches := EnvelopeRegex.FindStringSubmatch(input)
	if submatches != nil {
		return submatches[1]
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

func DecryptEnvelopes(input string, strategy Decryptor) (output string, err error) {
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

	output = EnvelopeRegex.ReplaceAllStringFunc(input, replfn)

	return

}

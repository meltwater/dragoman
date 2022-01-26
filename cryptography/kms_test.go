package cryptography

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type mockKmsClient struct {
	// GenerateDataKey response values
	GenerateDataKeyResponse *kms.GenerateDataKeyOutput
	GenerateDataKeyError    error

	// Decrypt response values
	DecryptResponse *kms.DecryptOutput
	DecryptError    error
}

func (c mockKmsClient) GenerateDataKey(ctx context.Context, input *kms.GenerateDataKeyInput, opts ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error) {
	return c.GenerateDataKeyResponse, c.GenerateDataKeyError
}

func (c mockKmsClient) Decrypt(ctx context.Context, input *kms.DecryptInput, opts ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	return c.DecryptResponse, c.DecryptError
}

func TestKmsEncryption(t *testing.T) {
	t.Run("it should encrypt the data", func(t *testing.T) {
		s := KmsCryptoStrategy{
			client: mockKmsClient{
				// TODO Fill in the response we want
			},
		}

		_ = s
		// TODO Test the encryption mocking out responses from kms
	})
}

func TestKmsDecryption(t *testing.T) {
	t.Run("it should decrypt the data", func(t *testing.T) {

	})
}

func TestKmsIntegration(t *testing.T) {
	t.Run("it should decrypt data it encrypts", func(t *testing.T) {

	})
}

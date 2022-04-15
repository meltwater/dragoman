package cryptography

import (
	"context"
	"fmt"
	"testing"

	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type smClientMock struct {
	mock.Mock
}

func (m *smClientMock) GetSecretValue(ctx context.Context, input *sm.GetSecretValueInput, opts ...func(*sm.Options)) (*sm.GetSecretValueOutput, error) {
	args := m.Called(ctx, input, opts)

	return args.Get(0).(*sm.GetSecretValueOutput), args.Error(1)
}

func getMockSecretsManagerStrategy() (strategy *SecretsManagerCryptoStrategy, smClient *smClientMock) {
	smClient = new(smClientMock)
	strategy = &SecretsManagerCryptoStrategy{
		client: smClient,
	}

	return
}

func generateMockSmEncryptedString(key string, secret string, output *string) {
	strategy, _ := getMockSecretsManagerStrategy()

	var encKey []byte = nil
	if key != "" {
		encKey = []byte(key)
	}

	*output, _ = strategy.Encrypt(encKey, secret)
}

func TestSecretsManagerCryptoStrategyBuilder(t *testing.T) {
	t.Run("it should generate the secrets manager client", func(t *testing.T) {
		s, err := NewSecretsManagerCryptoStrategy("us-east-1")

		assert.NotNil(t, s)
		assert.NotNil(t, s.client)
		assert.Nil(t, err)
	})
}

func TestSmEncrypt(t *testing.T) {
	t.Run("it should encrypt the secrets manager key", func(t *testing.T) {
		strategy, _ := getMockSecretsManagerStrategy()

		keyId := "aKey"

		encryptedKey, err := strategy.Encrypt(nil, keyId)

		assert.Nil(t, err)
		assert.NotNil(t, encryptedKey)
		assert.Contains(t, encryptedKey, "ENC[SECMAN,")
	})

	t.Run("it should encrypt the secrets manager key and the secret key if provided", func(t *testing.T) {
		strategy, _ := getMockSecretsManagerStrategy()

		keyId := "aKey"
		secretKey := "thekeyformyvalue"

		encryptedKey, err := strategy.Encrypt([]byte(secretKey), keyId)

		assert.Nil(t, err)
		assert.NotNil(t, encryptedKey)
		assert.Contains(t, encryptedKey, "ENC[SECMAN,")
	})
}

func TestSmDecrypt(t *testing.T) {
	t.Run("it should output the string returned by secrets manager", func(t *testing.T) {
		superSecret := "Jon Snow gets resurrected"
		var encrypted string
		generateMockSmEncryptedString("", "aKey", &encrypted)

		fmt.Println(encrypted)

		strategy, mockSm := getMockSecretsManagerStrategy()

		mockSm.On("GetSecretValue", context.TODO(), mock.Anything, mock.Anything).Return(
			&sm.GetSecretValueOutput{
				SecretString: &superSecret,
			}, nil)

		decrypted, err := strategy.Decrypt(encrypted)

		assert.Nil(t, err)
		assert.Equal(t, superSecret, string(decrypted))
	})

	t.Run("it should output the value in the JSON returned by SM for the key provided", func(t *testing.T) {
		superSecret := "{\"myKey\":\"Jon Snow gets resurrected\"}"
		var encrypted string
		generateMockSmEncryptedString("myKey", "aKey", &encrypted)

		fmt.Println(encrypted)

		strategy, mockSm := getMockSecretsManagerStrategy()

		mockSm.On("GetSecretValue", context.TODO(), mock.Anything, mock.Anything).Return(
			&sm.GetSecretValueOutput{
				SecretString: &superSecret,
			}, nil)

		decrypted, err := strategy.Decrypt(encrypted)

		assert.Nil(t, err)
		assert.Equal(t, "Jon Snow gets resurrected", string(decrypted))
	})
}

package cryptography

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type kmsClientMock struct {
	mock.Mock
}

func (m *kmsClientMock) GenerateDataKey(ctx context.Context, input *kms.GenerateDataKeyInput, opts ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error) {
	args := m.Called(ctx, input, opts)

	return args.Get(0).(*kms.GenerateDataKeyOutput), args.Error(1)
}

func (m *kmsClientMock) Decrypt(ctx context.Context, input *kms.DecryptInput, opts ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	args := m.Called(ctx, input, opts)

	return args.Get(0).(*kms.DecryptOutput), args.Error(1)
}

func getMockKmsStrategy() (strategy *KmsCryptoStrategy, kmsClient *kmsClientMock) {
	kmsClient = new(kmsClientMock)
	strategy = &KmsCryptoStrategy{
		client: kmsClient,
	}

	return
}

func TestKmsCryptoStrategyBuilder(t *testing.T) {
	t.Run("it should generate the kms client", func(t *testing.T) {
		strategy, err := NewKmsCryptoStrategy("us-east-1")

		assert.NotNil(t, strategy)
		assert.NotNil(t, strategy.client)
		assert.Nil(t, err)
	})
}

func TestKmsGenerateDataKey(t *testing.T) {
	t.Run("it should generate an encrypted key", func(t *testing.T) {
		strategy, mockKms := getMockKmsStrategy()

		keyId := "aKey"

		// Mock out the kms call
		gdkInput := &kms.GenerateDataKeyInput{
			KeyId:         &keyId,
			NumberOfBytes: &KMS_DATA_KEY_LENGTH,
		}

		gdkOutput := &kms.GenerateDataKeyOutput{
			Plaintext:      []byte("some plaintext that is 32 bytes "),
			CiphertextBlob: []byte("a CiphertextBlob"),
		}

		mockKms.On("GenerateDataKey", context.TODO(), gdkInput, mock.Anything).Return(gdkOutput, nil)

		// Test the Generate Data Key function
		dataKey, encryptedDataKey, err := strategy.GenerateDataKey(keyId)

		assert.Nil(t, err)
		assert.NotNil(t, dataKey)
		assert.Len(t, *dataKey, 32)
		assert.NotNil(t, encryptedDataKey)
		assert.Equal(t, gdkOutput.CiphertextBlob, encryptedDataKey)
		mockKms.AssertNumberOfCalls(t, "GenerateDataKey", 1)
	})

	t.Run("it should return an error if there is a kms failure", func(t *testing.T) {
		strategy, mockKms := getMockKmsStrategy()

		keyId := "aKey"

		// Mock out the kms call
		gdkInput := &kms.GenerateDataKeyInput{
			KeyId:         &keyId,
			NumberOfBytes: &KMS_DATA_KEY_LENGTH,
		}

		mockKms.On("GenerateDataKey", context.TODO(), gdkInput, mock.Anything).Return(&kms.GenerateDataKeyOutput{}, fmt.Errorf("An Error"))

		_, _, err := strategy.GenerateDataKey(keyId)

		assert.Error(t, err)
	})

	t.Run("it should return an error if there is an issue generating the NaCL key", func(t *testing.T) {
		strategy, mockKms := getMockKmsStrategy()

		keyId := "aKey"

		// Mock out the kms call
		gdkInput := &kms.GenerateDataKeyInput{
			KeyId:         &keyId,
			NumberOfBytes: &KMS_DATA_KEY_LENGTH,
		}

		gdkOutput := &kms.GenerateDataKeyOutput{
			Plaintext:      []byte("some other plaintext that is not 32 bytes"),
			CiphertextBlob: []byte("a CiphertextBlob"),
		}

		mockKms.On("GenerateDataKey", context.TODO(), gdkInput, mock.Anything).Return(gdkOutput, nil)

		_, _, err := strategy.GenerateDataKey(keyId)

		assert.Error(t, err)
	})
}

func TestKmsEncryption(t *testing.T) {
	t.Run("it should encrypt the data", func(t *testing.T) {
		strategy, mockKms := getMockKmsStrategy()

		keyId := "aKey"

		// Mock out the kms call
		gdkInput := &kms.GenerateDataKeyInput{
			KeyId:         &keyId,
			NumberOfBytes: &KMS_DATA_KEY_LENGTH,
		}

		gdkOutput := &kms.GenerateDataKeyOutput{
			Plaintext:      []byte("some plaintext that is 32 bytes "),
			CiphertextBlob: []byte("a CiphertextBlob"),
		}

		mockKms.On("GenerateDataKey", context.TODO(), gdkInput, mock.Anything).Return(gdkOutput, nil)

		superSecret := "Jon Snow is a Targaryen"
		encryptedString, err := strategy.Encrypt([]byte(superSecret), keyId)

		assert.Nil(t, err)
		assert.NotNil(t, encryptedString)
		assert.Contains(t, encryptedString, "ENC[KMS,")
		assert.NotContains(t, encryptedString, superSecret)
		mockKms.AssertNumberOfCalls(t, "GenerateDataKey", 1)
	})

	t.Run("it should return an error if GenerateDataKey fails", func(t *testing.T) {
		strategy, mockKms := getMockKmsStrategy()

		keyId := "aKey"

		// Mock out the kms call
		gdkInput := &kms.GenerateDataKeyInput{
			KeyId:         &keyId,
			NumberOfBytes: &KMS_DATA_KEY_LENGTH,
		}

		mockKms.On("GenerateDataKey", context.TODO(), gdkInput, mock.Anything).Return(&kms.GenerateDataKeyOutput{}, fmt.Errorf("An Error"))

		superSecret := "Jon Snow is a Targaryen"
		encryptedString, err := strategy.Encrypt([]byte(superSecret), keyId)

		assert.Error(t, err)
		assert.Equal(t, "", encryptedString)
	})
}

func generateMockEncryptedString(key string, secret string, output *string) {
	strategy, mockKms := getMockKmsStrategy()

	// Mock out the kms call
	gdkInput := &kms.GenerateDataKeyInput{
		KeyId:         &key,
		NumberOfBytes: &KMS_DATA_KEY_LENGTH,
	}

	gdkOutput := &kms.GenerateDataKeyOutput{
		Plaintext:      []byte("some plaintext that is 32 bytes "),
		CiphertextBlob: []byte("a CiphertextBlob"),
	}

	mockKms.On("GenerateDataKey", context.TODO(), gdkInput, mock.Anything).Return(gdkOutput, nil)

	*output, _ = strategy.Encrypt([]byte(secret), key)
}

func TestKmsDecryption(t *testing.T) {
	t.Run("it should decrypt the data encrypted by Encrypt", func(t *testing.T) {
		superSecret := "Jon Snow is a Targaryen"
		var encrypted string
		generateMockEncryptedString("aKey", superSecret, &encrypted)

		strategy, mockKms := getMockKmsStrategy()

		mockDecryptInput := &kms.DecryptInput{
			CiphertextBlob: []byte("a CiphertextBlob"),
		}

		mockDecryptOutput := &kms.DecryptOutput{
			Plaintext: []byte("some plaintext that is 32 bytes "),
		}

		mockKms.On("Decrypt", context.TODO(), mockDecryptInput, mock.Anything).Return(mockDecryptOutput, nil)

		decrypted, err := strategy.Decrypt(encrypted)

		assert.Nil(t, err)
		assert.Equal(t, superSecret, string(decrypted))
	})

	t.Run("it should return an error if the key decryption fails", func(t *testing.T) {
		superSecret := "Jon Snow is a Targaryen"
		var encrypted string
		generateMockEncryptedString("aKey", superSecret, &encrypted)

		strategy, mockKms := getMockKmsStrategy()

		mockDecryptInput := &kms.DecryptInput{
			CiphertextBlob: []byte("a CiphertextBlob"),
		}

		mockKms.On("Decrypt", context.TODO(), mockDecryptInput, mock.Anything).Return(&kms.DecryptOutput{}, fmt.Errorf("oopsie"))

		decrypted, err := strategy.Decrypt(encrypted)

		assert.Error(t, err)
		assert.Len(t, decrypted, 0)
	})
}

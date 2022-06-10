package cryptography

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	KMS_DATA_KEY_LENGTH int32  = 32
	CRYPTO_KEY_KMS      string = "KMS"
)

type kmsEnvelopeEncryptionPayload struct {
	EncryptedDataKey []byte
	Nonce            *[24]byte
	Message          []byte
}

// kmsCryproClientIfc allows us to mock the kms client in tests
type kmsCryptoClientIfc interface {
	GenerateDataKey(context.Context, *kms.GenerateDataKeyInput, ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error)
	Decrypt(context.Context, *kms.DecryptInput, ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

// KmsCryptoStrategy handles AWS KMS based encryption and decryption
type KmsCryptoStrategy struct {
	client kmsCryptoClientIfc
}

func NewKmsCryptoStrategy(region string) (*KmsCryptoStrategy, error) {
	// Do we need to set the region?
	var cfg aws.Config
	var err error

	if region == "" {
		if cfg, err = config.LoadDefaultConfig(context.TODO()); err != nil {
			return nil, err
		}
	} else {
		if cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region)); err != nil {
			return nil, err
		}
	}

	return &KmsCryptoStrategy{
		client: kms.NewFromConfig(cfg),
	}, nil
}

func (cs KmsCryptoStrategy) Key() string {
	return CRYPTO_KEY_KMS
}

func (cs *KmsCryptoStrategy) GenerateDataKey(keyId string) (*[32]byte, []byte, error) {
	// Use KMS to generate a data key
	var resp *kms.GenerateDataKeyOutput
	var err error

	if resp, err = cs.client.GenerateDataKey(context.TODO(), &kms.GenerateDataKeyInput{
		KeyId:         &keyId,
		NumberOfBytes: aws.Int32(KMS_DATA_KEY_LENGTH),
	}); err != nil {
		return nil, nil, err
	}

	// Convert from byte slice to byte array
	var dataKey *[32]byte
	if dataKey, err = AsNaCLKey(resp.Plaintext); err != nil {
		return nil, nil, err
	}

	return dataKey, resp.CiphertextBlob, nil
}

func (cs KmsCryptoStrategy) Encrypt(payload []byte, key string) (string, error) {
	var (
		dataKey          *[32]byte
		encryptedDataKey []byte
		err              error
	)

	// Use KMS to generate the data key
	if dataKey, encryptedDataKey, err = cs.GenerateDataKey(key); err != nil {
		return "", err
	}

	// Initialize the payload for the envelope
	envelopePayload := &kmsEnvelopeEncryptionPayload{
		EncryptedDataKey: encryptedDataKey,
		Nonce:            &[24]byte{},
	}

	// Generate the nonce
	if _, err = io.ReadFull(rand.Reader, envelopePayload.Nonce[:]); err != nil {
		return "", fmt.Errorf("failed to generate random nonce: %v", err)
	}

	// Seal the envelope
	envelopePayload.Message = secretbox.Seal(
		envelopePayload.Message,
		payload,
		envelopePayload.Nonce,
		dataKey)

	buff := &bytes.Buffer{}
	if err = gob.NewEncoder(buff).Encode(envelopePayload); err != nil {
		return "", err
	}

	return WrapEncoding("KMS", buff.Bytes()), nil
}

func (cs KmsCryptoStrategy) Decrypt(input string) ([]byte, error) {
	encrypted, err := UnwrapEncoding(input)
	if err != nil {
		return nil, fmt.Errorf("unable to unwrap the encrypted secret: %v", err)
	}

	// Decode the payload struct
	var payload kmsEnvelopeEncryptionPayload
	if err = gob.NewDecoder(bytes.NewReader(encrypted)).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode the message payload: %v", err)
	}

	// Decrypt the key
	var resp *kms.DecryptOutput
	if resp, err = cs.client.Decrypt(
		context.TODO(),
		&kms.DecryptInput{
			CiphertextBlob: payload.EncryptedDataKey,
		}); err != nil {
		return nil, fmt.Errorf("unable to decipher the kms key: %v", err)
	}

	// Convert the key to the expected NaCL type
	var key *[32]byte
	if key, err = AsNaCLKey(resp.Plaintext); err != nil {
		return nil, fmt.Errorf("unable to read kms key: %v", err)
	}

	// Decrypt the message
	var plaintext []byte
	var ok bool
	if plaintext, ok = secretbox.Open(plaintext, payload.Message, payload.Nonce, key); !ok {
		return nil, fmt.Errorf("failed to open the envelope")
	}

	return plaintext, nil
}

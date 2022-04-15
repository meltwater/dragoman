package cryptography

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type smEnvelopeEncryptionPayload struct {
	SecretID  []byte // Secret ARN or Name
	SecretKey []byte // Key for Secret Key/Value pairs
}

type smCryptoClientIfc interface {
	GetSecretValue(context.Context, *sm.GetSecretValueInput, ...func(*sm.Options)) (*sm.GetSecretValueOutput, error)
}

type SecretsManagerCryptoStrategy struct {
	client smCryptoClientIfc
}

func NewSecretsManagerCryptoStrategy(region string) (*SecretsManagerCryptoStrategy, error) {
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

	return &SecretsManagerCryptoStrategy{
		client: sm.NewFromConfig(cfg),
	}, nil
}

// Encrypt will generate the wrapped encoded string
// @param: payload is used for key/value pair keys
// @param: key is expected to be the key of the secret that will be used for decryption
// @returns: The base64 encoded arn with the encryption strategy key
func (cs SecretsManagerCryptoStrategy) Encrypt(payload []byte, key string) (string, error) {
	envelopePayload := &smEnvelopeEncryptionPayload{
		SecretID:  []byte(key),
		SecretKey: payload,
	}

	buff := &bytes.Buffer{}
	if err := gob.NewEncoder(buff).Encode(envelopePayload); err != nil {
		return "", err
	}

	return WrapEncoding("SECMAN", buff.Bytes()), nil
}

// Decrypt will pull the secret from Secrets Manager
func (cs SecretsManagerCryptoStrategy) Decrypt(input string) ([]byte, error) {
	// Unwrap the payload
	encrypted, err := UnwrapEncoding(input)
	if err != nil {
		return nil, fmt.Errorf("unable to unwrap the secret key: %v", err)
	}

	// Decode the payload struct
	var payload smEnvelopeEncryptionPayload
	if err = gob.NewDecoder(bytes.NewReader(encrypted)).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode the message payload: %v", err)
	}

	secretId := string(payload.SecretID)

	// Pull the key from KMS
	var resp *sm.GetSecretValueOutput
	if resp, err = cs.client.GetSecretValue(
		context.TODO(),
		&sm.GetSecretValueInput{
			SecretId: &secretId,
		}); err != nil {
		return nil, fmt.Errorf("unable to decipher the secret: %v", err)
	}

	secretString := *resp.SecretString
	if secretString == "" {
		return nil, fmt.Errorf("only string secrets are currently supported")
	}

	if payload.SecretKey != nil {
		secrets := map[string]string{}
		json.Unmarshal([]byte(secretString), &secrets)

		return []byte(secrets[string(payload.SecretKey)]), nil
	}

	return []byte(secretString), nil
}

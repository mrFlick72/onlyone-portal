package cypto

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/kms"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/awsclient"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
)

type AwsKmsKayRepository struct {
	kmsClient *kms.Client
	storage   map[string]string // keyId → base64-encoded KMS ciphertext blob
	cache     map[string][]byte // keyId → decripted key in memory cache
}

func NewAwsKmsKeyRepository(ctx context.Context) (KeyRepository, error) {
	cfg, err := awsclient.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("aws kms key repository: load config: %w", err)
	}

	configManager := config.GetConfigurationManagerInstance()
	keyId := configManager.GetConfigFor("key.aws-kms.storage.key")
	encryptedKeyValue := configManager.GetConfigFor("key.aws-kms.storage.key-value")

	return &AwsKmsKayRepository{
		kmsClient: kms.NewFromConfig(cfg),
		storage:   map[string]string{keyId: encryptedKeyValue},
		cache:     make(map[string][]byte, 0),
	}, nil
}

func (r *AwsKmsKayRepository) GetKeyFor(keyId string) (SymmetricKey, error) {
	if plainTextKey, ok := r.cache[keyId]; ok {
		return SymmetricKey{content: plainTextKey}, nil
	} else {
		encryptedValue, ok := r.storage[keyId]
		if !ok {
			return SymmetricKey{}, fmt.Errorf("key %q not found", keyId)
		}

		ciphertext, err := base64.StdEncoding.DecodeString(encryptedValue)
		if err != nil {
			return SymmetricKey{}, fmt.Errorf("decode kms ciphertext for key %q: %w", keyId, err)
		}

		result, err := r.kmsClient.Decrypt(context.Background(), &kms.DecryptInput{
			KeyId:          &keyId,
			CiphertextBlob: ciphertext,
		})

		if err != nil {
			return SymmetricKey{}, fmt.Errorf("kms decrypt key %q: %w", keyId, err)
		}
		r.cache[keyId] = result.Plaintext
		return SymmetricKey{content: result.Plaintext}, nil
	}

}

package cypto

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/go-playground/assert/v2"
)

func newAwsKmsKayRepositoryWithCache() *AwsKmsKayRepository {
	keyId := "A_KEY_ID"
	encryptedKeyValue := "AN_ENCRYTPED_KEY"

	return &AwsKmsKayRepository{
		kmsClient: newLocalKmsClient(),
		storage:   map[string]string{keyId: encryptedKeyValue},
		cache:     map[string][]byte{keyId: []byte(encryptedKeyValue)},
		mu:        sync.RWMutex{},
	}
}

func newAwsKmsKayRepository(storage map[string]string) *AwsKmsKayRepository {
	return &AwsKmsKayRepository{
		kmsClient: newLocalKmsClient(),
		storage:   storage,
		cache:     make(map[string][]byte, 0),
		mu:        sync.RWMutex{},
	}
}
func masterKeySetup() (string, []byte, []byte) {
	kmsClient := newLocalKmsClient()

	masterKey, err := kmsClient.CreateKey(context.TODO(), &kms.CreateKeyInput{
		KeyUsage: types.KeyUsageTypeEncryptDecrypt,
	})

	if err != nil {
		fmt.Printf("Error on key creation: %v", err)
		panic("AWS KMS key generation error")
	}

	key, err := kmsClient.GenerateDataKey(context.TODO(),
		&kms.GenerateDataKeyInput{
			KeyId:   masterKey.KeyMetadata.KeyId,
			KeySpec: types.DataKeySpecAes256,
		})

	if err != nil {
		fmt.Printf("Error on data key creation: %v", err)
		panic("AWS KMS data key generation error")
	}

	fmt.Println("key")
	fmt.Println(key)
	return *masterKey.KeyMetadata.KeyId, key.CiphertextBlob, key.Plaintext
}

func TestWhenAKeyIsInTheCache(t *testing.T) {
	uut := newAwsKmsKayRepositoryWithCache()
	actual, err := uut.GetKeyFor("A_KEY_ID")

	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("AN_ENCRYTPED_KEY"), actual.content)
}

func TestWhenAKeyIsNotInTheCache(t *testing.T) {
	masterKeyId, encryptedKey, plainTextKey := masterKeySetup()
	uut := newAwsKmsKayRepository(map[string]string{masterKeyId: base64.StdEncoding.EncodeToString(encryptedKey)})
	actual, _ := uut.GetKeyFor(masterKeyId)
	fmt.Println("actual")
	fmt.Println(actual)
	assert.Equal(t, plainTextKey, actual.content)
}

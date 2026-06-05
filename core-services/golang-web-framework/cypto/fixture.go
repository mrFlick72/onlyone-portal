//go:build test

package cypto

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

const localStackEndpoint = "http://localhost:4566"

func localStackConfig() (aws.Config, error) {
	return aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("xxx", "xxx", "xxx")),
		aws_config.WithRegion("eu-central-1"),
		aws_config.WithBaseEndpoint(localStackEndpoint),
	)
}

func newLocalKmsClient() *kms.Client {
	cfg, err := localStackConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return kms.NewFromConfig(cfg)
}


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

	return *masterKey.KeyMetadata.KeyId, key.CiphertextBlob, key.Plaintext
}

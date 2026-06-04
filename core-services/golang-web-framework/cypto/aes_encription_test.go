package cypto

import (
	"context"
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type inMemoryKeyRepository struct {
	keys map[string]SymmetricKey
}

func (r inMemoryKeyRepository) GetKeyFor(ctx context.Context, keyId string) (SymmetricKey, error) {
	key, ok := r.keys[keyId]
	if !ok {
		return SymmetricKey{}, fmt.Errorf("key %q not found", keyId)
	}

	return key, nil
}

func TestAesCbcCipherEncryptDecryptWithInMemoryKeyRepository(t *testing.T) {
	keyId := "test-key"
	repository := inMemoryKeyRepository{
		keys: map[string]SymmetricKey{
			keyId: {
				content: []byte("0123456789abcdef0123456789abcdef"),
			},
		},
	}
	cipher := AesCbcCipher{
		kyeId:      keyId,
		repository: repository,
	}
	plaintext := "secret message"

	encrypted, err := cipher.Encrypt(plaintext)

	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, plaintext, encrypted)

	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	require.NoError(t, err)
	assert.Zero(t, len(decoded)%aes.BlockSize)
	assert.Greater(t, len(decoded), aes.BlockSize)

	decrypted, err := cipher.Decrypt(encrypted)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

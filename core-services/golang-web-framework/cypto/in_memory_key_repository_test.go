package cypto

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemoryKeyRepositoryLoadsConfiguredKey(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "application.yml")
	err := os.WriteFile(configFile, []byte(`
key:
  in-memory:
    storage:
      key: configured-key
      key-value: 0123456789abcdef0123456789abcdef
`), 0600)
	require.NoError(t, err)
	t.Setenv("CONFIG_FILE_LOCATION", configFile)

	repository := NewInMemoryKeyRepository()

	key, err := repository.GetKeyFor(context.TODO(), "configured-key")
	require.NoError(t, err)
	assert.Equal(t, []byte("0123456789abcdef0123456789abcdef"), key.Content())
}

func TestInMemoryKeyRepositoryGetKeyFor(t *testing.T) {
	tests := []struct {
		name      string
		keyId     string
		wantKey   []byte
		assertErr assert.ErrorAssertionFunc
	}{
		{
			name:      "known key",
			keyId:     "test-key",
			wantKey:   []byte("0123456789abcdef0123456789abcdef"),
			assertErr: assert.NoError,
		},
		{
			name:      "unknown key",
			keyId:     "missing-key",
			assertErr: assert.Error,
		},
	}

	repository := &InMemoryKeyRepository{
		storage: map[string]string{
			"test-key": "0123456789abcdef0123456789abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := repository.GetKeyFor(context.TODO(), tt.keyId)

			tt.assertErr(t, err)
			assert.Equal(t, tt.wantKey, key.Content())
		})
	}
}

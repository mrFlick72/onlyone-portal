package cypto

import (
	"context"
	"fmt"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
)

type InMemoryKeyRepository struct {
	storage map[string]string
}

func NewInMemoryKeyRepository() KeyRepository {
	configManager := config.GetConfigurationManagerInstance()
	keyId := configManager.GetConfigFor("key.in-memory.storage.key")
	keyValue := configManager.GetConfigFor("key.in-memory.storage.key-value")
	return &InMemoryKeyRepository{storage: map[string]string{keyId: keyValue}}
}


func (receiver *InMemoryKeyRepository) GetKeyFor(ctx context.Context, keyId string) (SymmetricKey, error) {
	keyValue, ok := receiver.storage[keyId]
	if !ok {
		return SymmetricKey{}, fmt.Errorf("key %q not found", keyId)
	}

	return SymmetricKey{content: []byte(keyValue)}, nil
}

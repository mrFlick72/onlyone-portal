package cypto

import "github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"

type InMemoryKeyRepository struct {
	storage map[string]string
}

func NewInMemoryKeyRepository() KeyRepository {
	configManager := config.GetConfigurationManagerInstance()
	keyId := configManager.GetConfigFor("key.in-memory.storage.key")
	keyValue := configManager.GetConfigFor("key.in-memory.storage.keyValue")
	return &InMemoryKeyRepository{storage: map[string]string{keyId: keyValue}}
}

func (receiver *InMemoryKeyRepository) GetKeyFor(keyId string) (SymmetricKey, error) {
	return SymmetricKey{content: []byte(receiver.storage[keyId])}, nil
}

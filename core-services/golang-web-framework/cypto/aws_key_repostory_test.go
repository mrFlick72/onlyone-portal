package cypto

import (
	"context"
	"encoding/base64"
	"sync"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestWhenAKeyIsInTheCache(t *testing.T) {
	uut := newAwsKmsKayRepositoryWithCache()
	actual, err := uut.GetKeyFor(context.TODO(), "A_KEY_ID")

	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("AN_ENCRYTPED_KEY"), actual.content)
}

func TestWhenAKeyIsNotInTheCache(t *testing.T) {
	masterKeyId, encryptedKey, plainTextKey := masterKeySetup()
	uut := newAwsKmsKayRepository(map[string]string{masterKeyId: base64.StdEncoding.EncodeToString(encryptedKey)})
	actual, _ := uut.GetKeyFor(context.TODO(), masterKeyId)
	assert.Equal(t, plainTextKey, actual.content)
}

func TestWhenKeyIdIsNotFoundInStorage(t *testing.T) {
	uut := newAwsKmsKayRepository(map[string]string{})
	_, err := uut.GetKeyFor(context.TODO(), "MISSING_KEY")
	assert.NotEqual(t, nil, err)
}

func TestWhenStoredCiphertextIsInvalidBase64(t *testing.T) {
	uut := newAwsKmsKayRepository(map[string]string{"KEY": "not-base64!!!"})
	_, err := uut.GetKeyFor(context.TODO(), "KEY")
	assert.NotEqual(t, nil, err)
}

func TestCacheIsPopulatedAfterFirstCacheMiss(t *testing.T) {
	masterKeyId, encryptedKey, plainTextKey := masterKeySetup()
	uut := newAwsKmsKayRepository(map[string]string{masterKeyId: base64.StdEncoding.EncodeToString(encryptedKey)})

	first, err1 := uut.GetKeyFor(context.TODO(), masterKeyId)
	second, err2 := uut.GetKeyFor(context.TODO(), masterKeyId)

	assert.Equal(t, nil, err1)
	assert.Equal(t, nil, err2)
	assert.Equal(t, plainTextKey, first.content)
	assert.Equal(t, first.content, second.content)

	_, inCache := uut.cache[masterKeyId]
	assert.Equal(t, true, inCache)
}

func TestConcurrentAccessOnCachedKey(t *testing.T) {
	uut := newAwsKmsKayRepositoryWithCache()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key, err := uut.GetKeyFor(context.TODO(), "A_KEY_ID")
			assert.Equal(t, nil, err)
			assert.Equal(t, []byte("AN_ENCRYTPED_KEY"), key.content)
		}()
	}
	wg.Wait()
}

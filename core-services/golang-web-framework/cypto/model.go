package cypto

import "context"

type Key interface {
	Content() []byte
}

type SymmetricKey struct {
	content []byte
}

func (k SymmetricKey) Content() []byte {
	return k.content
}

type KeyRepository interface {
	GetKeyFor(ctx context.Context, keyId string) (Key, error)
}

type Cipher interface {
	Encrypt(ctx context.Context, plaintext string) (string, error)
	Decrypt(ctx context.Context, ciphertext string) (string, error)
}

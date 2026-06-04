package cypto

import "context"

type Key interface {
	Content() []byte
}

type SymmetricKey struct {
	content []byte
}

type KeyRepository interface {
	GetKeyFor(ctx context.Context, keyId string) (SymmetricKey, error)
}

type Cipher interface {
	Encrypt(ctx context.Context, plaintext string) (string, error)
	Decrypt(ctx context.Context, ciphertext string) (string, error)
}

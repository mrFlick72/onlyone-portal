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
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

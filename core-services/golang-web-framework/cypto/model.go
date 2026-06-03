package cypto

type Key interface {
	Content() []byte
}

type SymmetricKey struct {
	content []byte
}

type KeyRepository interface {
	GetKeyFor(keyId string) (SymmetricKey, error)
}

type Cipher interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

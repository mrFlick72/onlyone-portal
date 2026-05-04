package cypto

type Key struct {
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

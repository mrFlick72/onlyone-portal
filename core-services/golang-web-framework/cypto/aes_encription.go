package cypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type AesCbcCipher struct {
	kyeId      string
	repository KeyRepository
}

func (c *AesCbcCipher) Encrypt(plaintext string) (string, error) {
	key, err := c.repository.GetKeyFor(c.kyeId)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key.content)
	if err != nil {
		return "", err
	}

	// Pad plaintext to block size
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := append(([]byte(plaintext)), bytes.Repeat([]byte{byte(padding)}, padding)...)

	ciphertext := make([]byte, aes.BlockSize+len(padtext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padtext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *AesCbcCipher) Decrypt(ciphertext string) (string, error) {
	key, err := c.repository.GetKeyFor(c.kyeId)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key.content)
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) <= aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	if len(decodedCiphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decodedCiphertext, decodedCiphertext)

	// Unpad
	padding := int(decodedCiphertext[len(decodedCiphertext)-1])
	if padding == 0 || padding > aes.BlockSize || padding > len(decodedCiphertext) {
		return "", fmt.Errorf("invalid padding")
	}

	return string(decodedCiphertext[:len(decodedCiphertext)-padding]), nil
}

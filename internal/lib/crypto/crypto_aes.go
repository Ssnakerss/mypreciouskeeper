package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

// generateFixedKey uses SHA-256 to fix key length (32 байта).
func generateFixedKey(key []byte) []byte {
	hash := sha256.Sum256(key)
	return hash[:]
}

// encryptAES encrypr data with  key.
func EncryptAES(key []byte, plaintextBytes []byte) ([]byte, error) {
	fixedKey := generateFixedKey(key)

	block, err := aes.NewCipher(fixedKey)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]

	// Generate a random IV
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	return ciphertext, nil
}

// decryptAES decrypt data with  key.
func DecryptAES(key []byte, ciphertext []byte) ([]byte, error) {
	fixedKey := generateFixedKey(key)

	block, err := aes.NewCipher(fixedKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

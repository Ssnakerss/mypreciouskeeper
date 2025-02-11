package crypto

import "crypto/aes"

func EncryptAES(data []byte, key []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	res := make([]byte, aesblock.BlockSize())
	aesblock.Encrypt(res, data)
	return res, nil
}

func DecryptAES(data []byte, key []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	res := make([]byte, aesblock.BlockSize())
	aesblock.Decrypt(res, data)
	return res, nil
}

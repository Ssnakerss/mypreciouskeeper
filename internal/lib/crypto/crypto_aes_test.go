package crypto

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestEncryptAES(t *testing.T) {
	data := gofakeit.Sentence(10)
	t.Log(data)
	key := gofakeit.Sentence(1)
	t.Log(key)
	encrypted, err := EncryptAES(key, data)
	require.NoError(t, err)
	t.Log(encrypted)
	decrypted, err := DecryptAES(key, encrypted)
	require.NoError(t, err)
	t.Log(decrypted)

}

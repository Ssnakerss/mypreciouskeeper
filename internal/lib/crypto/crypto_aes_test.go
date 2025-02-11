package crypto

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestEncryptAES(t *testing.T) {
	data := gofakeit.Sentence(10)
	t.Log(data)
	key := "" //gofakeit.Sentence(1)
	t.Log(key)
	bkey := []byte(key)
	encrypted, err := EncryptAES(bkey, []byte(data))
	require.NoError(t, err)
	decrypted, err := DecryptAES(bkey, encrypted)
	require.NoError(t, err)
	t.Log(decrypted)
	require.Equal(t, data, string(decrypted))

}

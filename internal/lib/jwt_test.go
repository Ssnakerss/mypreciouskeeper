package lib

import (
	"testing"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestNewJWT(t *testing.T) {
	user := &models.User{
		ID:    1,
		Email: "test@test.com",
	}

	token, err := NewJWT(user, time.Hour)
	t.Log(token)
	require.NoError(t, err)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(AppSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.Equal(t, true, ok)

	t.Log(claims)

	tokenParsed, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("tututu"), nil
	})

	t.Log(err)
	require.Error(t, err)

	// email := claims["username"].(string)

	// type args struct {
	// 	user     *models.User
	// 	duration time.Duration
	// }
	// tests := []struct {
	// 	name            string
	// 	args            args
	// 	wantTokenString string
	// 	wantErr         bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		gotTokenString, err := NewJWT(tt.args.user, tt.args.duration)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("NewJWT() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		if gotTokenString != tt.wantTokenString {
	// 			t.Errorf("NewJWT() = %v, want %v", gotTokenString, tt.wantTokenString)
	// 		}
	// 	})
	// }
}

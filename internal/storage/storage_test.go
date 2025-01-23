package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/brianvoe/gofakeit"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestDBStorage_CreateUser(t *testing.T) {
	dsn := os.Getenv("POSTGRE_DSN")
	require.NotEmpty(t, dsn)

	db, err := New(context.Background(), dsn, time.Second*3)
	require.NoError(t, err)

	email := gofakeit.Email()
	// pass := gofakeit.Password(true, true, true, true, false, 10)

	//Testing user creation
	t.Log("testing create user")
	usr, err := db.CreateUser(context.Background(), email, "abc")
	require.NoError(t, err)

	//Testing user get
	t.Log("testting get user")
	usr, err = db.GetUser(context.Background(), email)
	require.NoError(t, err)
	require.NotEmpty(t, usr.Email)

	//Testing not existing  user get
	t.Log("testting not existing get user")
	usr, err = db.GetUser(context.Background(), "user_not-exist")
	require.NoError(t, err)
	require.Equal(t, int64(-1), usr.ID)

	//Testing duplicate user creation
	t.Log("testing duplicate user creation")
	usr, err = db.CreateUser(context.Background(), email, "abc")
	t.Log(err)
	require.Equal(t, apperrs.ErrUserAlreadyExists, err)
}

func TestDBStorage_saveGetAsset(t *testing.T) {
	dsn := os.Getenv("POSTGRE_DSN")
	require.NotEmpty(t, dsn)
	// dsn := "postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable"
	db, err := New(context.Background(), dsn, time.Second*3)
	require.NoError(t, err)

	str := `
	--------
	test
	test
	test
	--------
	`

	id, err := db.CreateAsset(context.Background(), 1, "text", "asset contains some text", []byte(str))
	if err != nil {
		t.Fatalf("save asset error: %v", err)
	}

	data, err := db.GetAsset(context.Background(), id)
	if err != nil {
		t.Fatalf("get asset error: %v", err)
	}

	t.Log(string(data), err)
}

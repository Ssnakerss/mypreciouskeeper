package storage

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func TestDBStorage_CreateUser(t *testing.T) {
	dsn := os.Getenv("POSTGRE_DSN")
	if dsn == "" {
		t.Fatal("dsn is not set")
		return
	}
	// dsn := "postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable"
	db, err := New(context.Background(), dsn, time.Second*3)
	if err != nil {
		t.Fatalf("db connection error: %v", err)
	}

	rnd := rand.New(rand.NewSource(uint64(time.Now().Nanosecond())))
	email := "email" + strconv.Itoa(rnd.Intn(256))

	//Testing user creation
	t.Log("testing create user")
	usr, err := db.CreateUser(context.Background(), email, "abc")
	if err != nil {
		t.Fatalf("user create error: %v", err)
	}

	//Testing user get
	t.Log("testting get user")
	usr, err = db.GetUser(context.Background(), email)
	if err != nil {
		t.Fatalf("user get error^ %v", err)
	}
	if usr.Email == "" {
		t.Fatalf("user get fail, email is empty, usr: %v", usr)
	}

	//Testing not existing  user get
	t.Log("testting not existing get user")
	usr, err = db.GetUser(context.Background(), "user_not-exist")
	if err != nil {
		t.Fatalf("user get error %v", err)
	}
	require.Equal(t, int64(-1), usr.ID)

	//Testing duplicate user creation
	t.Log("testtin duplicate user")
	usr, err = db.CreateUser(context.Background(), email, "abc")
	if err == nil {
		t.Fatalf("duplicate user create success, email: %v", email)
	}

}

func TestDBStorage_saveGetAsset(t *testing.T) {
	dsn := os.Getenv("POSTGRE_DSN")
	if dsn == "" {
		t.Fatal("dsn is not set")
		return
	}
	// dsn := "postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable"
	db, err := New(context.Background(), dsn, time.Second*3)
	if err != nil {
		t.Fatalf("db connection error: %v", err)
	}

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

package storage

import (
	"context"
	"strconv"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/exp/rand"
)

func TestDBStorage_CreateUser(t *testing.T) {
	// dsn := os.Getenv("POSTGRE_DSN")
	// if dsn == "" {
	// 	t.Fatal("dsn is not set")
	// 	return
	// }
	dsn := "postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable"
	db, err := New(context.Background(), dsn, time.Second*3)
	if err != nil {
		t.Fatalf("db connection error: %v", err)
	}

	rnd := rand.New(rand.NewSource(uint64(time.Now().Nanosecond())))
	email := "email" + strconv.Itoa(rnd.Intn(256))

	//Testing user creation
	t.Log("testtin create user")
	usr, err := db.CreateUser(context.Background(), email, "abc")
	if err != nil {
		t.Fatalf("user create error: %v", err)
	}

	//Testing user get
	t.Log("testtin get user")
	usr, err = db.GetUser(context.Background(), email)
	if err != nil {
		t.Fatalf("user get error^ %v", err)
	}
	if usr.Email == "" {
		t.Fatalf("user get fail, email is empty, usr: %v", usr)
	}

	//Testing duplicate user creation
	t.Log("testtin duplicate user")
	usr, err = db.CreateUser(context.Background(), email, "abc")
	if err == nil {
		t.Fatalf("duplicate user create success, email: %v", email)
	}

}

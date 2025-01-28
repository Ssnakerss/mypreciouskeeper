package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/domain/models"
	"github.com/brianvoe/gofakeit"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestDBStorage_User(t *testing.T) {
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

	//Testing duplicate user creation
	t.Log("testing duplicate user creation")
	usr, err = db.CreateUser(context.Background(), email, "abc")
	t.Log(err)
	require.Equal(t, apperrs.ErrUserAlreadyExists, err)
}

func TestDBStorage_Asset(t *testing.T) {
	dsn := os.Getenv("POSTGRE_DSN")
	require.NotEmpty(t, dsn)
	// dsn := "postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable"
	db, err := New(context.Background(), dsn, time.Second*3)
	require.NoError(t, err)

	str := gofakeit.Sentence(1)
	usrid := gofakeit.Int64()

	asset := &models.Asset{
		UserID:  usrid,
		Sticker: "test sticker here",
		Type:    "text",
		Body:    []byte(str),
	}

	t.Log("Create asset")
	id, err := db.CreateAsset(context.Background(), asset)
	require.NoError(t, err)

	t.Log("Get asset")
	rasset, err := db.GetAsset(context.Background(), usrid, id)
	require.NoError(t, err)
	require.Equal(t, asset, rasset)

	t.Log(asset)
	t.Log(rasset)

	t.Log("Update asset")
	rasset.Sticker = "updated sticker"
	err = db.UpdateAsset(context.Background(), rasset)
	require.NoError(t, err)
	asset, err = db.GetAsset(context.Background(), usrid, id)
	require.NoError(t, err)
	require.Equal(t, asset.Sticker, "updated sticker")

	t.Log("Delete asset")
	err = db.DeleteAsset(context.Background(), usrid, id)
	require.NoError(t, err)

	t.Log("Get deleted asset")
	rasset, err = db.GetAsset(context.Background(), usrid, id)
	require.Equal(t, apperrs.ErrAssetNotFound, err)

	t.Log("Update deleted asset")
	err = db.UpdateAsset(context.Background(), asset)
	require.Equal(t, apperrs.ErrAssetNotFound, err)

	t.Log("Testing ListAssets")
	//Create 3 asset, 2 same data, 1 with different sticker
	//Select by userid - count should be 3
	//Select by userid and type - count should be 2
	//Select by userid and sticker - count should be 1

	usrid = gofakeit.Int64()
	asset = &models.Asset{
		UserID:  usrid,
		Sticker: "test sticker here",
		Type:    "text",
		Body:    []byte(gofakeit.Sentence(10)),
	}
	_, err = db.CreateAsset(context.Background(), asset)
	require.NoError(t, err)
	_, err = db.CreateAsset(context.Background(), asset)
	require.NoError(t, err)
	asset = &models.Asset{
		UserID:  usrid,
		Sticker: "another sticker ",
		Type:    "card",
		Body:    []byte(gofakeit.Sentence(10)),
	}
	_, err = db.CreateAsset(context.Background(), asset)
	require.NoError(t, err)

	t.Log("Select by userid - count should be 3")
	assets, err := db.ListAssets(context.Background(), usrid, "", "")
	require.NoError(t, err)
	require.Equal(t, 3, len(assets))

	t.Log("Select by userid and type - count should be 2")
	assets, err = db.ListAssets(context.Background(), usrid, "text", "")
	require.NoError(t, err)
	require.Equal(t, 2, len(assets))

	t.Log("Select by userid and sticker - count should be 1")
	assets, err = db.ListAssets(context.Background(), usrid, "", "another")
	require.NoError(t, err)
	require.Equal(t, 1, len(assets))

}

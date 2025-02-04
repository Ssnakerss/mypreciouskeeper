package services

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/client/storage"
	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestServices_AuthSQLite(t *testing.T) {
	//Initialize SQLite db
	filePath := gofakeit.Sentence(1)
	require.NotEmpty(t, filePath)
	t.Log(filePath)

	db, err := storage.New(filePath+".db", time.Second*3)
	require.NoError(t, err)

	//Prepare tables
	errString := db.Prepare(context.Background())
	t.Log(errString)
	require.Empty(t, errString)

	//Logger setup
	l := logger.Setup("local")

	authService := NewAuthService(l, db, time.Hour*3)
	require.NotNil(t, authService)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, 10)

	t.Log("Testing user register")
	_, err = authService.RegisterUser(context.Background(), email, pass)
	require.NoError(t, err)

	t.Log("Testing existing user login")
	_, err = authService.Login(context.Background(), email, pass)
	require.NoError(t, err)

	//Fail cases
	t.Log("Testing same user register")
	_, err = authService.RegisterUser(context.Background(), email, pass)
	t.Log(err)
	require.Error(t, err)

	t.Log("Testing incorrect password login")
	_, err = authService.Login(context.Background(), email, "incorrect pass")
	t.Log(err)
	require.Error(t, err)

	t.Log("Testing incorrect user")
	_, err = authService.Login(context.Background(), "incorrect user", pass)
	t.Log(err)
	require.Error(t, err)

	//Cleanup

	err = db.Close()
	require.NoError(t, err)
	err = os.Remove(filePath + ".db")
	require.NoError(t, err)
}

func TestServices_AssetSQLite(t *testing.T) {
	//Initialize SQLite db
	filePath := gofakeit.Sentence(1)
	require.NotEmpty(t, filePath)
	t.Log(filePath)

	db, err := storage.New(filePath+".db", time.Second*3)
	require.NoError(t, err)

	//Prepare tables
	errString := db.Prepare(context.Background())
	t.Log(errString)
	require.Empty(t, errString)

	l := logger.Setup("local")
	assetService := NewAssetService(l, db)
	require.NotNil(t, assetService)
	authService := NewAuthService(l, db, time.Hour*3)
	require.NotNil(t, authService)

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, 10)

	t.Log("Testing user register")
	userID, err := authService.RegisterUser(context.Background(), email, pass)
	require.NoError(t, err)

	bodyStr := gofakeit.Sentence(1)
	aType := "TEXT"
	aSticker := gofakeit.Sentence(1)
	asset := &models.Asset{
		//Define asset ID in program logic
		ID:        gofakeit.Int64(),
		UserID:    userID,
		Sticker:   aSticker,
		Type:      aType,
		Body:      []byte(bodyStr),
		DeletedYN: "N",
	}

	t.Log("Create asset")
	asset, err = assetService.Create(context.Background(), asset)
	require.NoError(t, err)

	t.Log("Get asset")
	rasset, err := assetService.Get(context.Background(), userID, asset.ID)
	require.NoError(t, err)
	require.Equal(t, asset.ID, rasset.ID)

	t.Log("Update asset")
	rasset.Sticker = "updated sticker"
	err = assetService.Update(context.Background(), rasset)
	require.NoError(t, err)
	uasset, err := assetService.Get(context.Background(), userID, rasset.ID)
	require.NoError(t, err)
	require.Equal(t, uasset.Sticker, "updated sticker")

	t.Log("Get non-existing asset")
	_, err = assetService.Get(context.Background(), userID, 0)
	t.Log(err)
	require.Error(t, err)

	t.Log("Delete asset")
	err = assetService.Delete(context.Background(), userID, rasset.ID)
	require.NoError(t, err)

	//It is possible to get deleted asset locally
	// t.Log("Get deleted asset")

	//It is possible to updated delted asset locally
	// t.Log("Update deleted asset")

	t.Log("Testing ListAssets")
	//Create 3 asset, 2 same data, 1 with different sticker
	//Select by userid - count should be 3
	//Select by userid and type - count should be 2
	//Select by userid and sticker - count should be 1

	usrid := gofakeit.Int64()
	asset = &models.Asset{
		ID:        gofakeit.Int64(),
		UserID:    usrid,
		Sticker:   "test sticker here",
		Type:      models.AssetTypeMemo,
		Body:      []byte(gofakeit.Sentence(10)),
		DeletedYN: "N",
	}
	_, err = db.CreateAsset(context.Background(), asset)
	require.NoError(t, err)

	asset.ID = gofakeit.Int64()
	_, err = db.CreateAsset(context.Background(), asset)
	require.NoError(t, err)

	asset = &models.Asset{
		ID:        gofakeit.Int64(),
		UserID:    usrid,
		Sticker:   "another sticker ",
		Type:      models.AssetTypeCard,
		Body:      []byte(gofakeit.Sentence(10)),
		DeletedYN: "N",
	}
	_, err = db.CreateAsset(context.Background(), asset)
	require.NoError(t, err)

	t.Log("Select by userid - count should be 3")
	assets, err := db.ListAssets(context.Background(), usrid, "", "")
	require.NoError(t, err)
	require.Equal(t, 3, len(assets))

	t.Log("Select by userid and type - count should be 2")
	assets, err = db.ListAssets(context.Background(), usrid, models.AssetTypeMemo, "")
	require.NoError(t, err)
	require.Equal(t, 2, len(assets))

	t.Log("Select by userid and sticker - count should be 1")
	assets, err = db.ListAssets(context.Background(), usrid, "", "another")
	require.NoError(t, err)
	require.Equal(t, 1, len(assets))

	err = db.Close()
	require.NoError(t, err)
	err = os.Remove(filePath + ".db")
	require.NoError(t, err)
}

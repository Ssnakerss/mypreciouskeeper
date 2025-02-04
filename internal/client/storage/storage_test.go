package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/brianvoe/gofakeit"

	"github.com/stretchr/testify/require"
)

func TestDBStorage_User(t *testing.T) {
	filePath := gofakeit.Sentence(1)
	require.NotEmpty(t, filePath)
	t.Log(filePath)

	db, err := New(filePath+".db", time.Second*3)
	require.NoError(t, err)

	//Prepare tables
	errString := db.Prepare(context.Background())
	t.Log(errString)
	require.Empty(t, errString)

	//Testing user creation
	user := models.User{
		ID:        gofakeit.Int64(),
		Email:     gofakeit.Email(),
		PassHash:  gofakeit.Sentence(1),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Log("testing create user")
	usr, err := db.CreateUser(context.Background(), &user)
	require.NoError(t, err)

	//Testing user get
	t.Log("testting get user")
	usr, err = db.GetUser(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, usr.Email)

	//Testing not existing  user get
	t.Log("testting not existing get user")
	usr, err = db.GetUser(context.Background(), "user_not-exist")
	require.NoError(t, err)

	//Testing duplicate user creation
	t.Log("testing duplicate user creation")
	usr, err = db.CreateUser(context.Background(), &user)
	t.Log(err)
	require.Equal(t, apperrs.ErrUserAlreadyExists, err)

	err = db.Close()
	require.NoError(t, err)
	err = os.Remove(filePath + ".db")
	require.NoError(t, err)
}

func TestDBStorage_Asset(t *testing.T) {
	filePath := gofakeit.Sentence(1)
	require.NotEmpty(t, filePath)
	t.Log(filePath)

	db, err := New(filePath+".db", time.Second*3)
	require.NoError(t, err)

	//Prepare tables
	errString := db.Prepare(context.Background())
	t.Log(errString)
	require.Empty(t, errString)

	//Create test asset
	asset := &models.Asset{
		UserID:    gofakeit.Int64(),
		Sticker:   "test sticker here",
		Type:      "text",
		Body:      []byte(gofakeit.Sentence(100)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedYN: "N",
		// DeletedAt: time.Now(),
	}

	t.Log("Create asset")
	_, err = db.CreateAsset(context.Background(), asset)
	require.NoError(t, err)

	t.Log("Get asset")
	rasset, err := db.GetAsset(context.Background(), asset.UserID, asset.ID)

	require.NoError(t, err)
	require.Equal(t, asset.ID, rasset.ID)

	t.Log("Update asset")
	rasset.Sticker = "updated sticker"
	rasset.UpdatedAt = time.Now()

	err = db.UpdateAsset(context.Background(), rasset)
	require.NoError(t, err)

	asset, err = db.GetAsset(context.Background(), asset.UserID, asset.ID)
	require.NoError(t, err)
	require.Equal(t, asset.Sticker, "updated sticker")

	require.NotEqual(t, asset.UpdatedAt, rasset.UpdatedAt)

	// t.Log("Delete asset")
	// err = db.DeleteAsset(context.Background(), usrid, id)
	// require.NoError(t, err)

	// t.Log("Get deleted asset")
	// rasset, err = db.GetAsset(context.Background(), usrid, id)
	// require.Equal(t, apperrs.ErrAssetNotFound, err)

	// t.Log("Update deleted asset")
	// err = db.UpdateAsset(context.Background(), asset)
	// require.Equal(t, apperrs.ErrAssetNotFound, err)

	// t.Log("Testing ListAssets")
	// //Create 3 asset, 2 same data, 1 with different sticker
	// //Select by userid - count should be 3
	// //Select by userid and type - count should be 2
	// //Select by userid and sticker - count should be 1

	// usrid = gofakeit.Int64()
	// asset = &models.Asset{
	// 	UserID:  usrid,
	// 	Sticker: "test sticker here",
	// 	Type:    "text",
	// 	Body:    []byte(gofakeit.Sentence(10)),
	// }
	// _, err = db.CreateAsset(context.Background(), asset)
	// require.NoError(t, err)
	// _, err = db.CreateAsset(context.Background(), asset)
	// require.NoError(t, err)
	// asset = &models.Asset{
	// 	UserID:  usrid,
	// 	Sticker: "another sticker ",
	// 	Type:    "card",
	// 	Body:    []byte(gofakeit.Sentence(10)),
	// }
	// _, err = db.CreateAsset(context.Background(), asset)
	// require.NoError(t, err)

	// t.Log("Select by userid - count should be 3")
	// assets, err := db.ListAssets(context.Background(), usrid, "", "")
	// require.NoError(t, err)
	// require.Equal(t, 3, len(assets))

	// t.Log("Select by userid and type - count should be 2")
	// assets, err = db.ListAssets(context.Background(), usrid, "text", "")
	// require.NoError(t, err)
	// require.Equal(t, 2, len(assets))

	// t.Log("Select by userid and sticker - count should be 1")
	// assets, err = db.ListAssets(context.Background(), usrid, "", "another")
	// require.NoError(t, err)
	// require.Equal(t, 1, len(assets))

	err = db.Close()
	require.NoError(t, err)
	err = os.Remove(filePath + ".db")
	require.NoError(t, err)

}

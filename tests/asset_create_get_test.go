package tests

import (
	"encoding/json"
	"testing"

	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"
	"github.com/Ssnakerss/mypreciouskeeper/tests/suite"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Asset_Create(t *testing.T) {
	ctx, st := suite.New(t) // Создаём Suite

	// Generate fake email  and password for test
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, 10)

	// Creating client and make Register and Login request
	respReg, err := st.AClient.Register(ctx, &grpcserver.RegisterRequest{
		Email: email,
		Pass:  pass,
	})
	t.Log("user id ", respReg.GetUserId())
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AClient.Login(ctx, &grpcserver.LoginRequest{
		Email: email,
		Pass:  pass,
	})
	require.NoError(t, err)
	t.Log(respLogin)

	//Checking received token
	token := respLogin.GetToken()
	require.NotEmpty(t, token)
	t.Log(token)

	t.Log("Test creating asset")
	memo := models.Memo{
		Text: gofakeit.Sentence(100),
	}
	body, err := json.Marshal(memo)
	require.NoError(t, err)

	respAssetCreate, err := st.AssetClient.Create(ctx, &grpcserver.CreateRequest{
		Token:   token,
		Type:    "MEMO",
		Sticker: "ITEM FROM GRPC TEST",
		Body:    body,
	})
	require.NoError(t, err)
	t.Log(respAssetCreate)

	t.Log("Test get asset")
	respAssetGet, err := st.AssetClient.Get(ctx, &grpcserver.GetRequest{
		Token:   token,
		AssetId: respAssetCreate.AssetId,
	})
	require.NoError(t, err)
	t.Log(respAssetGet)

	t.Log("Test batch asset create and get")

	for i := 0; i < 10; i++ {
		memo := models.Memo{
			Text: gofakeit.Sentence(100),
		}
		body, err := json.Marshal(memo)
		require.NoError(t, err)
		respAssetCreate, err := st.AssetClient.Create(ctx, &grpcserver.CreateRequest{
			Token:   token,
			Type:    "MEMO",
			Sticker: "ITEM FROM GRPC TEST" + gofakeit.Sentence(1),
			Body:    body,
		})
		require.NoError(t, err)
		t.Log(respAssetCreate)
	}

	respList, err := st.AssetClient.List(ctx, &grpcserver.ListRequest{Token: token})
	require.NoError(t, err)
	t.Log(respList)

}

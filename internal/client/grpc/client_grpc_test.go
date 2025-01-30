package grpcClient

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func Test_gRPCClient_Register(t *testing.T) {
	cl := NewGRPCClient(net.JoinHostPort("localhost", "44044"))
	type args struct {
		email string
		pass  string
	}
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, false, false, 10)

	tests := []struct {
		name    string
		c       *GRPCClient
		args    args
		wantErr bool
	}{
		{
			name: "Normal Register",
			c:    cl,
			args: args{
				email: email,
				pass:  pass,
			},
			wantErr: false,
		},
		{
			name: "Empty pass Register",
			c:    cl,
			args: args{
				email: email,
				pass:  "",
			},
			wantErr: true,
		},
		{
			name: "Empty email Register",
			c:    cl,
			args: args{
				email: "",
				pass:  pass,
			},
			wantErr: true,
		},
		{
			name: "Empty pass and email Register",
			c:    cl,
			args: args{
				email: "",
				pass:  "",
			},
			wantErr: true,
		},
		{
			name: "Duplicate Register",
			c:    cl,
			args: args{
				email: email,
				pass:  pass,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.c.Register(tt.args.email, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("gRPCClient.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

// Test different Login cases
func Test_gRPCClient_Login(t *testing.T) {
	cl := NewGRPCClient(net.JoinHostPort("localhost", "44044"))
	type args struct {
		email string
		pass  string
	}
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, false, false, 10)
	//Register first
	_, err := cl.Register(email, pass)
	require.NoError(t, err)

	tests := []struct {
		name    string
		c       *GRPCClient
		args    args
		wantErr bool
	}{
		{
			name: "Normal Login",
			c:    cl,
			args: args{
				email: email,
				pass:  pass,
			},
			wantErr: false,
		},
		{
			name: "Empty pass Login",
			c:    cl,
			args: args{
				email: email,
				pass:  "",
			},
			wantErr: true,
		},
		{
			name: "Empty email Login",
			c:    cl,
			args: args{
				email: "",
				pass:  pass,
			},
			wantErr: true,
		},
		{
			name: "Empty pass and email Login",
			c:    cl,
			args: args{
				email: "",
				pass:  "",
			},
			wantErr: true,
		},
		{
			name: "Incorrect email login",
			c:    cl,
			args: args{
				email: "incorrect email",
				pass:  pass,
			},
			wantErr: true,
		},
		{
			name: "Incorrect pass login",
			c:    cl,
			args: args{
				email: email,
				pass:  "incorrect pass",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.c.Login(tt.args.email, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("gRPCClient.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

// Test Create and Get asset
func Test_gRPCClient_AssetCreateGet(t *testing.T) {
	cl := NewGRPCClient(net.JoinHostPort("localhost", "44044"))
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, false, false, 10)
	//Register first
	userid, err := cl.Register(email, pass)
	require.NoError(t, err)
	t.Log("userid: ", userid)
	//Login
	_, err = cl.Login(email, pass)
	require.NoError(t, err)

	card := models.Card{
		Name:     "Test Card",
		Number:   "1234567890123456",
		CVV:      "123",
		ExpMonth: 12,
		ExpYear:  2022,
	}
	body, err := json.Marshal(card)
	require.NoError(t, err)
	asset := &models.Asset{
		Type:    "CARD",
		Sticker: "test sample card",
		Body:    body,
	}
	//Creating asset
	assetId, err := cl.CreateAsset(asset)
	require.NoError(t, err)
	t.Log("asset id: ", assetId)

	//Getting asset by id
	getAsset, err := cl.GetAsset(assetId)
	require.NoError(t, err)
	t.Log(getAsset)

	//Comparing created and get asset

	rCard := models.Card{}
	err = json.Unmarshal(getAsset.Body, &rCard)
	require.NoError(t, err)
	require.Equal(t, card, rCard)
}

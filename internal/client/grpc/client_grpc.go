package grpcClient

import (
	"context"
	"log"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	AuthClient  grpcserver.AuthClient
	AssetClient grpcserver.AssetClient
	token       string
	Conn        *grpc.ClientConn
}

// NewGRPCClient create client with Auth and Asset Endpoints from gRPC server
func NewGRPCClient(grpcAddress string) *GRPCClient {
	//TO-DO: server address from  config
	// grpcAddress := net.JoinHostPort("localhost", "44044")

	Conn, err := grpc.NewClient(
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc server connection failed: %v", err)
	}
	authClient := grpcserver.NewAuthClient(Conn)
	assetClient := grpcserver.NewAssetClient(Conn)

	return &GRPCClient{
		AuthClient:  authClient,
		AssetClient: assetClient,
	}
}

// Login to remote server with email and password and receive auth token
func (c *GRPCClient) Login(email string, pass string) (token string, err error) {
	loginResp, err := c.AuthClient.Login(context.Background(), &grpcserver.LoginRequest{
		Email: email,
		Pass:  pass,
	})
	if err != nil {
		return "", err
	}
	c.token = loginResp.Token
	return loginResp.Token, nil
}

// Register to remote server with email and password and receive userid
func (c *GRPCClient) Register(email string, pass string) (userid int64, err error) {
	registerResp, err := c.AuthClient.Register(context.Background(), &grpcserver.RegisterRequest{
		Email: email,
		Pass:  pass,
	})

	if err != nil {
		return -1, err
	}
	return registerResp.UserId, nil
}

// CreateAsset send Create requestAssetRequest to remote server with auth token and receive asset id
func (c *GRPCClient) CreateAsset(asset *models.Asset) (assetId int64, err error) {
	if c.token == "" {
		return -1, apperrs.ErrEmptyToken
	}
	createAssetResp, err := c.AssetClient.Create(context.Background(), &grpcserver.CreateRequest{
		Token:   c.token,
		Type:    asset.Type,
		Sticker: asset.Sticker,
		Body:    asset.Body,
	})
	if err != nil {
		return -1, err
	}
	return createAssetResp.AssetId, nil
}

// GetAsset send Get requestAssetRequest to remote server with auth token and receive asset
func (c *GRPCClient) GetAsset(assetId int64) (asset *models.Asset, err error) {
	if c.token == "" {
		return nil, apperrs.ErrEmptyToken
	}
	getAssetResp, err := c.AssetClient.Get(context.Background(), &grpcserver.GetRequest{Token: c.token, AssetId: assetId})
	if err != nil {
		return nil, err
	}
	return &models.Asset{
		ID:        getAssetResp.AssetId,
		Type:      getAssetResp.Type,
		Sticker:   getAssetResp.Sticker,
		Body:      getAssetResp.Body,
		CreatedAt: time.Unix(getAssetResp.CreatedAt, 0),
		UpdatedAt: time.Unix(getAssetResp.UpdatedAt, 0),
	}, nil
}

func (c *GRPCClient) List(assetType string) (assets []*models.Asset, err error) {
	assetList, err := c.AssetClient.List(context.Background(), &grpcserver.ListRequest{
		Token: c.token,
		Type:  assetType,
	})
	if err != nil {
		return nil, err
	}
	for _, asset := range assetList.Assets {
		assets = append(assets, &models.Asset{
			ID:        asset.AssetId,
			Type:      asset.Type,
			Sticker:   asset.Sticker,
			Body:      asset.Body,
			CreatedAt: time.Unix(asset.CreatedAt, 0),
			UpdatedAt: time.Unix(asset.UpdatedAt, 0),
		})
	}
	return assets, nil
}

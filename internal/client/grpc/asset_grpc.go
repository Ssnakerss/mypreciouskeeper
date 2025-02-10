package grpcClient

import (
	"context"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"
)

// CreateAsset send Create requestAssetRequest to remote server with auth token and receive asset id
func (c *GRPCClient) Create(
	ctx context.Context,
	asset *models.Asset,
) (*models.Asset, error) {
	if c.token == "" {
		return asset, apperrs.ErrEmptyToken
	}

	createAssetResp, err := c.AssetClient.Create(
		context.Background(),
		&grpcserver.CreateRequest{
			Token:   c.token,
			Type:    asset.Type,
			Sticker: asset.Sticker,
			Body:    asset.Body,
		},
	)
	if err == nil {
		asset.ID = createAssetResp.AssetId
	}
	return asset, err
}

// GetAsset send Get requestAssetRequest to remote server with auth token and receive asset
func (c *GRPCClient) Get(
	ctx context.Context,
	userID int64,
	assetId int64,
) (*models.Asset, error) {
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

func (c *GRPCClient) List(
	ctx context.Context,
	userID int64,
	assetType string,
	assetSticker string,
) (assets []*models.Asset, err error) {
	assetList, err := c.AssetClient.List(context.Background(), &grpcserver.ListRequest{
		Token:   c.token,
		Type:    assetType,
		Sticker: assetSticker,
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

func (c *GRPCClient) Delete(
	ctx context.Context,
	userID int64,
	aid int64) error {
	_, err := c.AssetClient.Delete(context.Background(), &grpcserver.DeleteRequest{Token: c.token, AssetId: aid})
	return err
}

func (c *GRPCClient) Update(ctx context.Context, asset *models.Asset) error {
	_, err := c.AssetClient.Update(context.Background(), &grpcserver.UpdateRequest{
		Token:   c.token,
		AssetId: asset.ID,
		Type:    asset.Type,
		Sticker: asset.Sticker,
		Body:    asset.Body,
	})
	return err
}

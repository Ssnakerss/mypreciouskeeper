package client

import (
	"context"

	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
)

// Asset functions
func (app *ClientApp) CreateAsset(
	ctx context.Context,
	asset *models.Asset) (*models.Asset, error) {
	return app.remoteAssetService.Create(ctx, asset)
}
func (app *ClientApp) GetAsset(
	ctx context.Context,
	id int64) (*models.Asset, error) {
	//TODO: check for user id
	return app.remoteAssetService.Get(ctx, -1, id)
}
func (app *ClientApp) List(
	ctx context.Context,
	assetType string) ([]*models.Asset, error) {
	//TODO: check for user id and asset sticker
	return app.remoteAssetService.List(ctx, -1, assetType, "")
}

func (app *ClientApp) UpdateAsset(
	ctx context.Context,
	asset *models.Asset) error {
	return app.remoteAssetService.Update(ctx, asset)
}

package client

import (
	"context"

	"github.com/Ssnakerss/mypreciouskeeper/internal/lib/crypto"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
)

// Asset functions
// Add body data encryption with masterpassword
func (app *ClientApp) CreateAsset(
	ctx context.Context,
	asset *models.Asset) (*models.Asset, error) {
	var err error

	b := asset.Body
	//Encrypting body data with masterpassword before save to storage
	asset.Body, err = crypto.EncryptAES(app.GetMasterPass(), asset.Body)
	if err != nil {
		return asset, err
	}

	a, err := app.remoteAssetService.Create(ctx, asset)
	//cancel encryption befor return, just in case
	a.Body = b
	return a, err
}

// GetAsset from storage by id
// Add body decryption with masterpassword
func (app *ClientApp) GetAsset(
	ctx context.Context,
	id int64) (*models.Asset, error) {
	//TODO: check for user id

	asset, err := app.remoteAssetService.Get(ctx, -1, id)
	if err != nil {
		return asset, err
	}
	//decrypt asset body before return
	asset.Body, err = crypto.DecryptAES(app.GetMasterPass(), asset.Body)
	if err != nil {
		return asset, err
	}
	return asset, nil
}

// List assets from storage (without body)
func (app *ClientApp) List(
	ctx context.Context,
	assetType string) ([]*models.Asset, error) {
	//TODO: check for user id and asset sticker
	return app.remoteAssetService.List(ctx, -1, assetType, "")
}

// Update asset on storage
// Ass body encryption with masterpassword
func (app *ClientApp) UpdateAsset(
	ctx context.Context,
	asset *models.Asset) error {
	var err error
	//Encrypting body data with masterpassword before save to storage
	b := asset.Body
	asset.Body, err = crypto.EncryptAES(app.GetMasterPass(), asset.Body)
	if err != nil {
		//cancel encryption befor return
		asset.Body = b
		return err
	}
	return app.remoteAssetService.Update(ctx, asset)
}

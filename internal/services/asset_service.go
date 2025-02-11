package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
)

type AssetStorage interface {
	CreateAsset(ctx context.Context, asset *models.Asset) (*models.Asset, error)
	GetAsset(ctx context.Context, userid int64, aid int64) (*models.Asset, error)
	ListAssets(ctx context.Context, userid int64, atype string, asticker string) ([]*models.Asset, error)

	UpdateAsset(ctx context.Context, asset *models.Asset) error
	DeleteAsset(ctx context.Context, userid int64, aid int64) error

	Close() error
}

type AssetService struct {
	l *slog.Logger
	s AssetStorage
}

func NewAssetService(l *slog.Logger, s AssetStorage) *AssetService {
	return &AssetService{
		l: l,
		s: s,
	}
}

// Create new asser record in sto
// gRPC mapping -  Create
func (a *AssetService) Create(
	ctx context.Context,
	asset *models.Asset,
) (*models.Asset, error) {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "AssetService.Create"
	l := a.l.With(slog.String("who", who),
		slog.String("type", asset.Type),
		slog.String("sticker", asset.Sticker),
	)
	l.Info("registering new asset")

	return a.s.CreateAsset(ctx, asset)
}

// Get asset data from storage
// gRPC mapping - Get
func (a *AssetService) Get(
	ctx context.Context,
	userID int64,
	aid int64,
) (*models.Asset, error) {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "AssetService.Get"
	l := a.l.With(slog.String("who", who),
		slog.Int64("id", aid),
	)
	l.Info("getting asset data by id")
	return a.s.GetAsset(ctx, userID, aid)
}

func (a *AssetService) List(
	ctx context.Context,
	userID int64,
	atype string,
	asticker string,
) ([]*models.Asset, error) {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "AssetService.List"
	l := a.l.With(slog.String("who", who),
		slog.Int64("user_id", userID),
		slog.String("type", atype),
		slog.String("sticker", asticker),
	)
	l.Info("getting asset data by id")
	return a.s.ListAssets(ctx, userID, atype, asticker)
}

// ListLatest for sync service - return recently  updated record
// Dummy to meet SyncService interface requirements
// TODO - implement
func (c *AssetService) ListLatest(
	ctx context.Context,
	userID int64,
	lastUpdated time.Time,
) (assets []*models.Asset, err error) {
	return nil, errors.New("not implemented")
}

// Update to latest asset data in storage
// Now - update all new asset
// TODO  -  check local update/create time
func (c *AssetService) UpdateToLatest(ctx context.Context,
	asset *models.Asset,
) (int, error) {
	err := c.Update(ctx, asset)
	if errors.Is(err, apperrs.ErrAssetNotFound) {
		//Update failed, tryung to create
		_, err = c.Create(ctx, asset)
	}
	return 0, err
}

// Update asset information in storage
// gRPC mapping - Update
func (a *AssetService) Update(
	ctx context.Context,
	asset *models.Asset,
) error {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "AssetService.Update"
	l := a.l.With(slog.String("who", who),
		slog.Int64("id", asset.ID),
		slog.String("type", asset.Type),
		slog.String("sticker", asset.Sticker),
	)
	l.Info("updating asset data by id")
	return a.s.UpdateAsset(ctx, asset)

}

// Delete asset data from storage by ID
// gRPC mapping - Delete
func (a *AssetService) Delete(
	ctx context.Context,
	userID int64,
	aid int64) error {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "AssetService.Delete"
	l := a.l.With(slog.String("who", who),
		slog.Int64("id", aid),
	)
	l.Info("deleting asset data by id")
	return a.s.DeleteAsset(ctx, userID, aid)
}

// Close closes underlying storage
func (a *AssetService) Close() {
	a.s.Close()
}

package services

import (
	"context"
	"log/slog"

	"github.com/Ssnakerss/mypreciouskeeper/internal/domain/models"
)

type AssetStorage interface {
	CreateAsset(ctx context.Context, asset *models.Asset) (int64, error)
	GetAsset(ctx context.Context, userid int64, aid int64) (*models.Asset, error)

	ListAssets(ctx context.Context, userid int64, atype string, asticker string) ([]*models.Asset, error)

	UpdateAsset(ctx context.Context, asset *models.Asset) error
	DeleteAsset(ctx context.Context, userid int64, aid int64) error
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
	userID int64,
	atype string,
	asticker string,
	abody []byte,
) (ad int64, err error) {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "AssetService.Create"
	l := a.l.With(slog.String("who", who),
		slog.String("type", atype),
		slog.String("sticker", asticker),
	)
	l.Info("registering new asset")
	newAsset := &models.Asset{UserID: userID, Type: atype, Sticker: asticker, Body: abody}

	return a.s.CreateAsset(ctx, newAsset)
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

// Update asset information in storage
// gRPC mapping - Update
func (a *AssetService) Update(
	ctx context.Context,
	userID int64,
	aid int64,
	atype string,
	asticker string,
	abody []byte,
) error {
	//who - current function name
	//for logging purpose to identify which function is calling
	who := "AssetService.Update"
	l := a.l.With(slog.String("who", who),
		slog.Int64("id", aid),
		slog.String("type", atype),
		slog.String("sticker", asticker),
	)
	l.Info("updating asset data by id")
	return a.s.UpdateAsset(ctx, &models.Asset{UserID: userID, ID: aid, Type: atype, Sticker: asticker, Body: abody})

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

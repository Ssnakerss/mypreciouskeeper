package server

import (
	"context"

	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AssetService interface {
	Create(
		ctx context.Context,
		asset *models.Asset,
	) (*models.Asset, error)

	Update(
		ctx context.Context,
		asset *models.Asset,
	) error

	Get(
		ctx context.Context,
		userID int64,
		aid int64,
	) (*models.Asset, error)

	List(
		ctx context.Context,
		userID int64,
		atype string,
		asticker string,
	) ([]*models.Asset, error)

	Delete(
		ctx context.Context,
		userID int64,
		aid int64,
	) error
}

type serverAssetAPI struct {
	grpcserver.UnimplementedAssetServer
	assetService AssetService
}

func NewServerAssetAPI(a AssetService) *serverAssetAPI {
	return &serverAssetAPI{
		assetService: a,
	}
}

func (s *serverAssetAPI) RegisterGRPC(gRPCServer *grpc.Server) {
	grpcserver.RegisterAssetServer(gRPCServer, s)
}

// Crate asset record in storage and return asset id
func (s *serverAssetAPI) Create(ctx context.Context, req *grpcserver.CreateRequest) (*grpcserver.CreateResponse, error) {
	//Verify authorization token
	user, err := verifyJWTPayload(req.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}
	asset := &models.Asset{
		UserID:  user.ID,
		Type:    req.Type,
		Sticker: req.Sticker,
		Body:    req.Body,
	}

	//Create asset
	aid, err := s.assetService.Create(ctx, asset)
	if err != nil {
		return nil, err
	}

	return &grpcserver.CreateResponse{AssetId: aid.ID}, err
}

// Get retrive asset by id and user id from storage
func (s *serverAssetAPI) Get(ctx context.Context, req *grpcserver.GetRequest) (*grpcserver.GetResponse, error) {
	//Verify authorization token
	user, err := verifyJWTPayload(req.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	asset, err := s.assetService.Get(ctx, user.ID, req.AssetId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Asset not found")
	}
	return &grpcserver.GetResponse{
		AssetId:   asset.ID,
		Type:      asset.Type,
		Sticker:   asset.Sticker,
		Body:      asset.Body,
		CreatedAt: asset.CreatedAt.Unix(),
		UpdatedAt: asset.UpdatedAt.Unix(),
	}, nil
}

// List return assets from storage selected by user id, type and sticker
func (s *serverAssetAPI) List(ctx context.Context, req *grpcserver.ListRequest) (*grpcserver.ListResponse, error) {
	//Verify authorization token
	user, err := verifyJWTPayload(req.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}
	assets, err := s.assetService.List(ctx, user.ID, req.Type, req.Sticker)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Asset not found")
	}
	rassets := []*grpcserver.ListResponse_Asset{}
	for _, a := range assets {
		rassets = append(rassets, &grpcserver.ListResponse_Asset{
			AssetId:   a.ID,
			Type:      a.Type,
			Sticker:   a.Sticker,
			Body:      a.Body,
			CreatedAt: a.CreatedAt.Unix(),
			UpdatedAt: a.UpdatedAt.Unix(),
		})

	}
	return &grpcserver.ListResponse{
		Assets: rassets,
	}, nil
}

// Update asset in storage by asset id and user id with new body and type and sticker
func (s *serverAssetAPI) Update(ctx context.Context, req *grpcserver.UpdateRequest) (*grpcserver.UpdateResponse, error) {
	//Verify authorization token
	user, err := verifyJWTPayload(req.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}
	asset := &models.Asset{
		ID:      req.AssetId,
		UserID:  user.ID,
		Type:    req.Type,
		Sticker: req.Sticker,
		Body:    req.Body,
	}
	err = s.assetService.Update(ctx, asset)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Asset not found")
	}
	return &grpcserver.UpdateResponse{
		AssetId: req.AssetId,
	}, nil

}

// Delete asset from storage by asset id and user id
func (s *serverAssetAPI) Delete(ctx context.Context, req *grpcserver.DeleteRequest) (*grpcserver.DeleteResponse, error) {
	//Verify authorization token
	user, err := verifyJWTPayload(req.GetToken())

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}
	err = s.assetService.Delete(ctx, user.ID, req.AssetId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Asset not found")

	}
	return &grpcserver.DeleteResponse{
		AssetId: req.AssetId,
	}, nil
}

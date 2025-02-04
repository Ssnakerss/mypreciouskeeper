package client

import (
	grpcClient "github.com/Ssnakerss/mypreciouskeeper/internal/client/grpc"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
)

// Service which provides authorization and user registration
type AuthService interface {
	Login(login, password string) (string, error)
	Register(login, password string) (int64, error)
	//TO-DO: implement Logout method
	// Logout() error
}

type AssetService interface {
	CreateAsset(asset *models.Asset) (int64, error)
	GetAsset(id int64) (*models.Asset, error)
	List(assetType string) ([]*models.Asset, error)
	//To-DO: implement UpdateAsset method and DeleteAsset methods
	// UpdateAsset(id int64, name, description string) (int64, error)
	// DeleteAsset(id int64) error
}

var App *ClientApp

type ClientApp struct {
	AuthService  AuthService
	AssetService AssetService
	AuthToken    string
	UserID       int64
	UserName     string
}

func NewClientApp(
	gRPCAddress string,
) *ClientApp {

	//TO-DO:  switch to common service
	grpc := grpcClient.NewGRPCClient(gRPCAddress)

	return &ClientApp{
		AuthService:  grpc,
		AssetService: grpc,
	}
}

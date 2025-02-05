package client

import (
	"context"
	"log"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/client/config"
	grpcClient "github.com/Ssnakerss/mypreciouskeeper/internal/client/grpc"
	"github.com/Ssnakerss/mypreciouskeeper/internal/client/storage"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/Ssnakerss/mypreciouskeeper/internal/services"
)

const (
	LOCAL  = "LOCAL"
	REMOTE = "REMOTE"
)

// Service which provides authorization and user registration
type AuthService interface {
	Login(
		ctx context.Context,
		login, password string) (string, error)
	Register(
		ctx context.Context,
		login, password string) (int64, error)
	//TODO: implement Logout method
	// Logout() error
}

type AssetService interface {
	Create(
		ctx context.Context,
		asset *models.Asset,
	) (*models.Asset, error)
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
	Update(
		ctx context.Context,
		asset *models.Asset,
	) error
	Delete(
		ctx context.Context,
		userID int64,
		aid int64) error
}

var App *ClientApp

type ClientApp struct {
	//Auth and asset services -  both remote and local
	remoteAuthService  AuthService
	remoteAssetService AssetService
	localAuthService   AuthService
	localAssetService  AssetService

	Workmode string //Workmode - LOCAL  or REMOTE  (LOCAL - local user, REMOTE - remote user)

	AuthToken    string //JWT token for gRPC remote Endpoints
	UserID       int64  //UserID when remote Login used
	LocalUsersID int64  //UserID when local Login used
	UserName     string

	//Keep login and password for remote auto login
	//if first login was offline - after connection established try login remotely
	login    string
	password string

	l   *slog.Logger
	cfg *config.Config

	//TODO: implement screens logic ???   how to add tea screeen into app
	// appScreens map[string]tea.Model
}

// NewClientApp create main app struct, initialize and assign local and remote data providers
func NewClientApp(
	l *slog.Logger,
	cfg *config.Config,
) *ClientApp {

	//Prepare local storage and service
	db, err := storage.New(cfg.StoragePath, time.Second*3)
	if err != nil {
		log.Fatal(err)
	}
	//Prepare tables
	db.Prepare(context.Background())

	authService := services.NewAuthService(l, db, time.Hour*3)
	assetService := services.NewAssetService(l, db)

	//Prepare remote gRPC service
	grpc := grpcClient.NewGRPCClient(net.JoinHostPort(cfg.GRPC.Host, strconv.Itoa(cfg.GRPC.Port)))

	return &ClientApp{
		l:                  l,
		cfg:                cfg,
		remoteAuthService:  grpc,
		remoteAssetService: grpc,
		localAuthService:   authService,
		localAssetService:  assetService,
	}
}

// Auth functions

// Login using remote service first
// If success - register user locally
// If fail - try to login locally
// If fail - return error
func (app *ClientApp) Login(
	ctx context.Context,
	login, password string) (string, error) {
	app.login = login
	app.password = password
	//Try login remotely
	token, err := app.remoteAuthService.Login(ctx, login, password)
	if err != nil {
		//Try login locally
		token, err = app.localAuthService.Login(ctx, login, password)
		if err != nil {
			//Cannot login
			return "", err
		}
		//Local login success
		//return empty  token
		app.Workmode = LOCAL
		return "", nil
	} else {
		//Remote login success
		app.Workmode = REMOTE
		//Try register same user locally with same login and password
		app.LocalUsersID, err = app.localAuthService.Register(ctx, login, password)
	}
	app.AuthToken = token
	return token, nil
}

// Register user with remote service first
// Then register with local service
func (app *ClientApp) Register(
	ctx context.Context,
	login, password string) (int64, error) {
	//Try register remotely
	app.l.Info("trying remote regsiter")
	remoteUserID, err := app.remoteAuthService.Register(ctx, login, password)
	if err != nil {
		//Try register locally
		app.l.Error("remote regsiter", "error", err)
		app.l.Info("trying local regsiter")
		localUserID, err := app.localAuthService.Register(ctx, login, password)
		if err != nil {
			app.l.Error("local regsiter", "error", err)
			return 0, err
		}
		//Local register success
		return localUserID, nil
	}
	return remoteUserID, nil
}

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

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
	//TODO: implement Update method
	//Update()

	Close()
}

type AssetService interface {
	//Create asset
	Create(
		ctx context.Context,
		asset *models.Asset,
	) (*models.Asset, error)
	//Get asset by id and userid
	Get(
		ctx context.Context,
		userID int64,
		aid int64,
	) (*models.Asset, error)
	//List all assets by userid and with like Type and Sticker
	List(
		ctx context.Context,
		userID int64,
		atype string,
		asticker string,
	) ([]*models.Asset, error)

	//Update asset
	Update(
		ctx context.Context,
		asset *models.Asset,
	) error
	Delete(
		ctx context.Context,
		userID int64,
		aid int64) error

	Close()
}

type PingService interface {
	Ping(ctx context.Context) (int64, error)
}

var App *ClientApp

type ClientApp struct {
	//Auth and asset services -  both remote and local
	remoteAuthService  AuthService
	remoteAssetService AssetService

	localAuthService  AuthService
	localAssetService AssetService

	pingService PingService
	SyncClient  *SyncClient

	Workmode string //Workmode - LOCAL  or REMOTE  (LOCAL - local user, REMOTE - remote user)

	AuthToken    string //JWT token for gRPC remote Endpoints
	RemoteUserID int64  //UserID when remote Login used
	LocalUsersID int64  //UserID when local Login used
	UserName     string

	//Keep login and password for remote auto login
	//if first login was offline - after connection established try login remotely
	login    string
	password string

	L   *slog.Logger
	cfg *config.Config

	Version   string
	BuildTime string

	//Context to sync goroutines
	SyncCtx       context.Context
	SyncCtxCancel context.CancelFunc

	//TODO: implement screens logic ???   how to add tea screeen into app
	// appScreens map[string]tea.Model
}

// NewClientApp create main app struct, initialize and assign local and remote data providers
func NewClientApp(
	ctx context.Context,
	l *slog.Logger,
	cfg *config.Config,
	v string, //version
	b string, //build time
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
	myGrpc := grpcClient.NewGRPCClient(net.JoinHostPort(cfg.GRPC.Host, strconv.Itoa(cfg.GRPC.Port)))

	//Setup sync
	syncClient := NewSyncClient(ctx, myGrpc, assetService)

	return &ClientApp{
		L:                  l,
		cfg:                cfg,
		remoteAuthService:  myGrpc,
		remoteAssetService: myGrpc,
		pingService:        myGrpc,

		localAuthService:  authService,
		localAssetService: assetService,

		SyncClient: syncClient,

		Version:   v,
		BuildTime: b,
	}
}

func (c *ClientApp) Ping(
	ctx context.Context,
	syncCtxCancel context.CancelFunc,
	interval int,
) {
	l := c.L.With("who", "ClientApp.Ping")

	pingTimeTicker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			l.Info("ping process terminated")
			syncCtxCancel()
			return
		case <-pingTimeTicker.C:
			if i, err := c.pingService.Ping(ctx); err != nil {
				l.Info("gRPC connection is not ready")
				c.Workmode = LOCAL
			} else {
				c.Workmode = REMOTE
				remoteTime := time.Unix(i, 0)
				l.Info("gRPC connection is ready", "remote time", remoteTime)

				//Restore remote login if user initially logged in locally
				if c.AuthToken == "" && c.login != "" && c.password != "" {
					l.Info("trying to login remotely")
					if token, err := c.remoteAuthService.Login(ctx, c.login, c.password); err != nil {
						l.Error("remote login failed", "err", err)
					} else {
						c.AuthToken = token
						c.RemoteUserID = i
						l.Info("remote login success")
					}
				}
				//Starting Sync process
				if c.RemoteUserID > 0 {
					updatedrecord, err := c.SyncClient.Sync(ctx, c.RemoteUserID)
					if err != nil {
						l.Error("Sync failed", "err", err)
					} else {
						l.Info("Sync success", "updated records", updatedrecord)
					}
				}
			}
		}
	}
}

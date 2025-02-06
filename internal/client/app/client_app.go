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

	Workmode string //Workmode - LOCAL  or REMOTE  (LOCAL - local user, REMOTE - remote user)

	AuthToken    string //JWT token for gRPC remote Endpoints
	RemoteUserID int64  //UserID when remote Login used
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
	ctx context.Context,
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
	myGrpc := grpcClient.NewGRPCClient(net.JoinHostPort(cfg.GRPC.Host, strconv.Itoa(cfg.GRPC.Port)))

	return &ClientApp{
		l:                  l,
		cfg:                cfg,
		remoteAuthService:  myGrpc,
		remoteAssetService: myGrpc,
		pingService:        myGrpc,

		localAuthService:  authService,
		localAssetService: assetService,
	}
}

func (c *ClientApp) Ping(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.l.Info("ping process terminated")
			return
		default:
			if i, err := c.pingService.Ping(ctx); err != nil {
				c.l.Info("gRPC connection is not ready")
				c.Workmode = LOCAL
			} else {
				c.Workmode = REMOTE
				remoteTime := time.Unix(i, 0)
				c.l.Info("gRPC connection is ready", "remote time", remoteTime)
			}
		}
		time.Sleep(time.Second * 10)
	}
}

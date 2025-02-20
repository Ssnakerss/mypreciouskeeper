package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/client/config"
	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
)

// Run application, initialize main App struct  and start services
func Run(
	baseCtx context.Context,
	version string,
	buildTime string,
) {

	cfg := config.Load()
	filename := fmt.Sprintf("app%s.log", time.Now().Format("20060102"))

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	l := logger.Setup(cfg.Env, file)
	l = l.With("who", "server/main")
	l.Info("server starting ...")

	App = NewClientApp(baseCtx, l, cfg, version, buildTime)

	App.SyncCtx, App.SyncCtxCancel = context.WithCancel(context.Background())
	defer App.SyncCtxCancel()

	//Start ping for remote service
	App.L.Info("starting ping seriver")
	go App.Ping(baseCtx, App.SyncCtxCancel, 5)

	//Initialize and start Sync service

	App.L.Info("app started")
}

// Stop application and close all services
func (c *ClientApp) Stop() {
	<-c.SyncCtx.Done()

	c.remoteAssetService.Close()
	c.remoteAssetService.Close()

	c.localAssetService.Close()
	c.localAuthService.Close()

	c.L.Info("app exit.")
}

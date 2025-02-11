package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
	server "github.com/Ssnakerss/mypreciouskeeper/internal/server/app"
	"github.com/Ssnakerss/mypreciouskeeper/internal/server/config"
)

func main() {

	cfg := config.Load()

	l := logger.Setup(cfg.Env, os.Stdout)
	l = l.With("who", "server/main")
	l.Info("server starting ...")

	app := server.New(l, cfg)

	go app.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Waiting for SIGINT (pkill -2) or SIGTERM
	<-stop

	// initiate graceful shutdown
	app.Shutdown()
}

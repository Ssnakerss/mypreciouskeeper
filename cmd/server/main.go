package main

import (
	"github.com/Ssnakerss/mypreciouskeeper/internal/config"
	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
)

func main() {

	//TODO: config initialize
	cfg := config.Load()

	//TODO: logger initialize
	l := logger.Setup(cfg.Env)
	l = l.With("who", "server/main")
	l.Info("server starting ...")

	//TODO: app initialize

	//TODO: gRPC start
}

package main

import (
	"github.com/Ssnakerss/mypreciouskeeper/internal/config"
	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
)

func main() {

	//TODO: config initialize
	cfg := config.Load()

	//TODO: logger initialize
	log := logger.Setup(cfg.Env)
	log = log.With("who", "server/main")
	log.Info("server starting ...")

	//TODO: app initialize

	//TODO: gRPC start
}

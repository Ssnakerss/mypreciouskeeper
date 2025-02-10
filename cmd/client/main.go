package main

import (
	"context"
	"log"

	"os"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/Ssnakerss/mypreciouskeeper/internal/client/config"
	"github.com/Ssnakerss/mypreciouskeeper/internal/client/screens"
	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func main() {
	cfg := config.Load()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	l := logger.Setup(cfg.Env, file)
	l = l.With("who", "server/main")
	l.Info("server starting ...")

	//Base app context
	baseCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client.App = client.NewClientApp(baseCtx, l, cfg)
	//Start ping for remote service
	syncCtx, syncCancel := context.WithCancel(context.Background())
	defer syncCancel()
	go client.App.Ping(baseCtx, syncCancel, 5)

	//Setup  initial app screen
	initialScreen := screens.RootScreen()

	//Clear screen before display initial screen
	output := termenv.NewOutput(os.Stdout)
	output.ClearScreen()

	//Start app with Tea screens
	//syscal signals are handled by  Tea program
	if _, err := tea.NewProgram(
		initialScreen,
	).Run(); err != nil {

		l.Error("Error running program:", logger.Err(err))
	}
	//TODO:....
	cancel()
	<-syncCtx.Done()
	client.App.Close()
	l.Info("app exit.")
}

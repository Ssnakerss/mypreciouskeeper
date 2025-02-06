package main

import (
	"context"
	"fmt"
	"log"

	"os"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/Ssnakerss/mypreciouskeeper/internal/client/config"
	"github.com/Ssnakerss/mypreciouskeeper/internal/client/screens"
	"github.com/Ssnakerss/mypreciouskeeper/internal/lib"
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
	defer func() { l.Info("program terminated!!!!!!") }()

	//TODO: add pre-shutdown tasks - db close etc
	go lib.SysCallProcess(baseCtx, cancel, l)

	client.App = client.NewClientApp(baseCtx, l, cfg)

	//Start ping for remote service
	go client.App.Ping(baseCtx)

	//Setup  initial app screen
	initialScreen := screens.RootScreen()

	//Clear screen before display initial screen
	output := termenv.NewOutput(os.Stdout)
	output.ClearScreen()

	//Start app with Tea screens
	if _, err := tea.NewProgram(
		initialScreen,
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

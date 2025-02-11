package main

import (
	"context"

	"os"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/Ssnakerss/mypreciouskeeper/internal/client/screens"
	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

// Var for build version and build time setup by ldflags
var (
	Version   string
	BuildTime string
)

func main() {

	//Base app context
	baseCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//run app
	client.Run(baseCtx, Version, BuildTime)

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
		client.App.L.Error("Error running program:", logger.Err(err))
	}

	client.App.L.Info("prepare to stop")
	cancel()
	client.App.Stop()
	client.App.L.Info("bye-bye")

}

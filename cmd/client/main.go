package main

import (
	"fmt"
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

	client.App = client.NewClientApp(l, cfg)
	//Setup  initial app screen
	initialScreen := screens.RootScreen()

	//Clear screen before display initial screen
	output := termenv.NewOutput(os.Stdout)
	output.ClearScreen()

	if _, err := tea.NewProgram(
		initialScreen,
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

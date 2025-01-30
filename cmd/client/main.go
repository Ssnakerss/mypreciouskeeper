package main

import (
	"fmt"
	"net"

	"os"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/Ssnakerss/mypreciouskeeper/internal/client/screens"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func main() {
	client.App = client.NewClientApp(net.JoinHostPort("localhost", "44044"))
	initialScreen := screens.RootScreen() // screens.ListMethodsScreen()

	output := termenv.NewOutput(os.Stdout)
	output.ClearScreen()

	//TO-DO - appy config params
	//Initial global singleton for App

	if _, err := tea.NewProgram(initialScreen).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

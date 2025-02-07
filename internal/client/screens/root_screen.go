package screens

import (
	"fmt"
	"strings"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	tea "github.com/charmbracelet/bubbletea"
)

type rootScreenModel struct {
	model tea.Model // this will hold the current screen model
}

func RootScreen() rootScreenModel {
	var rootModel tea.Model

	if client.App.UserName == "" {
		rootModel = AuthMenuScreen()
	} else {
		rootModel = CreateActionsMenuScreen()
	}

	return rootScreenModel{
		model: rootModel,
	}
}

// Wrapper for current screen Init View and Update methods
func (m rootScreenModel) Init() tea.Cmd {
	return m.model.Init()
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m rootScreenModel) View() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s", m.model.View())

	//Connection status 'widget'
	statusWidget(client.App.Workmode, &b)

	//Return view
	return b.String()
}

// SwitchScreen is the switcher which will switch between screens
func (m rootScreenModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.model = model
	//Return .Init() to initialize the screen
	return m.model, m.model.Init()
}

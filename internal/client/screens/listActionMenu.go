package screens

import (
	"os"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func ActionMenuScreen() ActionMenuModel {
	items := []list.Item{
		item{title: "Create", alias: "create"},
		item{title: "List", alias: "list"},
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "MY PRECIOUS KEEPER"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	output := termenv.NewOutput(os.Stdout)
	return ActionMenuModel{list: l, output: output}
}

type ActionMenuModel struct {
	list     list.Model
	quitting bool
	output   *termenv.Output
}

func (m ActionMenuModel) Init() tea.Cmd {
	return tea.SetWindowTitle("M_P_K")
}

func (m ActionMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			m.output.ClearScreen()

			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			m.quitting = true
			m.output.ClearScreen()

			return m, tea.Quit

		case tea.KeyEnter:
			item, ok := m.list.SelectedItem().(item)

			if ok {
				if item.alias == "register" {
					screen_y := ScreenRegister()
					return RootScreen().SwitchScreen(&screen_y)
				}

				if item.alias == "login" {
					screen_y := ScreenLogin()
					return RootScreen().SwitchScreen(&screen_y)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ActionMenuModel) View() string {
	return "\n" + m.list.View() + "\n" + "User ID: " + client.App.UserName
}

package screens

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func AuthMenuScreen() AuthMenuModel {
	items := []list.Item{
		item{title: "Register", alias: "register"},
		item{title: "Login", alias: "login"},
	}

	l := list.New(items, itemDelegate{}, defaultWidth, defaultHeight)
	l.Title = "MY PRECIOUS KEEPER"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	output := termenv.NewOutput(os.Stdout)
	return AuthMenuModel{list: l, output: output}
}

type AuthMenuModel struct {
	list     list.Model
	quitting bool
	output   *termenv.Output
}

func (m AuthMenuModel) Init() tea.Cmd {
	return tea.SetWindowTitle("M_P_K")
}

func (m AuthMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
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

func (m AuthMenuModel) View() string {
	return "\n" + m.list.View()
}

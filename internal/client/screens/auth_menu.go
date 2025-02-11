package screens

import (
	"fmt"
	"os"
	"strings"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("206"))

	return AuthMenuModel{
		list:    l,
		output:  output,
		spinner: sp,
	}
}

type AuthMenuModel struct {
	list     list.Model
	quitting bool
	output   *termenv.Output
	spinner  spinner.Model
}

func (m AuthMenuModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
	)
}

func (m AuthMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

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
					return RootScreen().SwitchScreen(&screen_y, "")
				}
				if item.alias == "login" {
					screen_y := ScreenLogin()
					return RootScreen().SwitchScreen(&screen_y, "")
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m AuthMenuModel) View() string {
	var b strings.Builder
	b.WriteString(m.list.View())

	fmt.Fprintf(&b, "\nVersion: %s | Bulid time: %s \n", client.App.Version, client.App.BuildTime)

	//Connection status 'widget'
	statusWidget(client.App.Workmode, &b)
	b.WriteString(m.spinner.View())

	return b.String()
}

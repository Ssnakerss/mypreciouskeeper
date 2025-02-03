package screens

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func CreateActionsMenuScreen() CreateActionsMenuModel {
	items := []list.Item{
		item{title: "Memo", alias: "create_memo"},
		item{title: "Credentials", alias: "create_cred"},
		item{title: "Card", alias: "create_card"},
		item{title: "File", alias: "create_file"},
		item{title: "View assets", alias: "view_assets"},
	}

	l := list.New(items, itemDelegate{}, defaultWidth, defaultHeight)
	l.Title = "MY PRECIOUS KEEPER"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	output := termenv.NewOutput(os.Stdout)
	return CreateActionsMenuModel{list: l, output: output}
}

type CreateActionsMenuModel struct {
	list     list.Model
	quitting bool
	output   *termenv.Output
}

func (m CreateActionsMenuModel) Init() tea.Cmd {
	return tea.SetWindowTitle("M_P_K")
}

func (m CreateActionsMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		//Quit the program when the user presses Ctrl-C
		case tea.KeyCtrlC:
			m.quitting = true
			m.output.ClearScreen()
			return m, tea.Quit
			//Return to prev screen
		case tea.KeyEsc, tea.KeyCtrlQ:
			screen_y := AuthMenuScreen()
			return RootScreen().SwitchScreen(&screen_y)
		case tea.KeyEnter:
			item, ok := m.list.SelectedItem().(item)

			if ok {
				if item.alias == "create_memo" {
					screen_y := CreateMemoScreen()
					return RootScreen().SwitchScreen(&screen_y)
				}
				if item.alias == "create_cred" {
					screen_y := CreateCredentialsScreen()
					return RootScreen().SwitchScreen(&screen_y)
				}
				if item.alias == "create_card" {
					// screen_y := CreateMemoScreen()
					// return RootScreen().SwitchScreen(&screen_y)
				}
				if item.alias == "create_file" {
					// screen_y := CreateMemoScreen()
					// return RootScreen().SwitchScreen(&screen_y)
				}
				if item.alias == "view_assets" {
					screen_y := CreateListScreen()
					return RootScreen().SwitchScreen(&screen_y)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m CreateActionsMenuModel) View() string {
	return "\n" + m.list.View()
}

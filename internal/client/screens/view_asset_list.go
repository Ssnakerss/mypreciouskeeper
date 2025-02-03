package screens

import (
	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type listItem struct {
	models.Asset
}

func (i listItem) Title() string       { return i.Type }
func (i listItem) Description() string { return i.Sticker }
func (i listItem) FilterValue() string { return i.Type }

type ListView struct {
	list list.Model
}

var footer string

func CreateListScreen() ListView {
	assetList, _ := client.App.GRPC.List("")
	Items := []list.Item{}
	for _, asset := range assetList {
		Items = append(Items,
			listItem{
				models.Asset{
					ID:      asset.ID,
					Type:    asset.Type,
					Sticker: asset.Sticker,
					Body:    asset.Body,
				},
			})
	}

	m := ListView{list: list.New(Items, list.NewDefaultDelegate(), defaultWidth, defaultHeight)}
	m.list.Title = "My Precious"
	return m
}

func (m ListView) Init() tea.Cmd {
	return nil
}

func (m ListView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			footer = m.list.SelectedItem().(listItem).Type
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			screen_y := CreateActionsMenuScreen()
			return RootScreen().SwitchScreen(&screen_y)
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListView) View() string {
	return docStyle.Render(m.list.View()) + "\n\n" + docStyle.Render(footer)
}

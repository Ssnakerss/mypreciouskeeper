package screens

import (
	"context"

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
	//REceive assets from the server to display
	assetList, _ := client.App.List(context.Background(), "")
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
			//Selected item - need to display details
			//Check type and display appropriate screen
			switch m.list.SelectedItem().(listItem).Type {
			case models.AssetTypeMemo:
				screen_y := CreateMemoScreen(m.list.SelectedItem().(listItem).ID)
				return RootScreen().SwitchScreen(&screen_y, models.AssetTypeMemo)
			case models.AssetTypeCredentials:
				screen_y := CreateCredentialsScreen(m.list.SelectedItem().(listItem).ID)
				return RootScreen().SwitchScreen(&screen_y, models.AssetTypeCredentials)
			case models.AssetTypeCard:
				screen_y := CreateCardScreen(m.list.SelectedItem().(listItem).ID)
				return RootScreen().SwitchScreen(&screen_y, models.AssetTypeCard)
			case models.AssetTypeFile:
				screen_y := CreateFileScreen(m.list.SelectedItem().(listItem).ID)
				return RootScreen().SwitchScreen(&screen_y, models.AssetTypeFile)
			default:
				footer = "NOT IMPLEMENTED:" + m.list.SelectedItem().(listItem).Type
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			screen_y := CreateActionsMenuScreen()
			return RootScreen().SwitchScreen(&screen_y, "")
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

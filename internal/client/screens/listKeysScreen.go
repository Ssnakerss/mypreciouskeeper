package screens

import (
	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type listItem struct {
	title, desc string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }

type ListView struct {
	list list.Model
}

var footer string

func CreateList() List {

	assetList, err := client.App.GRPC.List("")
	if err != nil {
		return
	}

	Items := []list.Item{
		listItem{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		listItem{title: "Nutella", desc: "It's good on toast"},
		listItem{title: "Bitter melon", desc: "It cools you down"},
		listItem{title: "Nice socks", desc: "And by that I mean socks without holes"},
		listItem{title: "Eight hours of sleep", desc: "I had this once"},
		listItem{title: "Cats", desc: "Usually"},
		listItem{title: "Plantasia, the album", desc: "My plants love it too"},
		listItem{title: "Pour over coffee", desc: "It takes forever to make though"},
		listItem{title: "VR", desc: "Virtual reality...what is there to say?"},
		listItem{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		listItem{title: "Linux", desc: "Pretty much the best OS"},
		listItem{title: "Business school", desc: "Just kidding"},
		listItem{title: "Pottery", desc: "Wet clay is a great feeling"},
		listItem{title: "Shampoo", desc: "Nothing like clean hair"},
		listItem{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		listItem{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		listItem{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		listItem{title: "Stickers", desc: "The thicker the vinyl the better"},
		listItem{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		listItem{title: "Warm light", desc: "Like around 2700 Kelvin"},
		listItem{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		listItem{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		listItem{title: "Terrycloth", desc: "In other words, towel fabric"},
	}

	m := ListView{list: list.New(Items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"
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
			footer = m.list.SelectedItem().(listItem).title

		case tea.KeyCtrlC:
			return m, tea.Quit
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

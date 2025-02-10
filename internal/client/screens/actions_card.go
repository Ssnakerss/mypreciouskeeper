package screens

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"context"
	"encoding/json"
	"fmt"

	"os"
	"strings"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type cardScreen struct {
	focusIndex int
	textInputs []textinput.Model
	cursorMode cursor.Mode

	caption string
	action  string //for button caption

	err     error
	success string

	asset *models.Asset
}

func CreateCardScreen(assetID int64) cardScreen {
	m := cardScreen{
		//6 fields to input
		//Sticker + Name, number, expire month/year,CVV
		textInputs: make([]textinput.Model, 6),
	}
	var err error
	//Create new asset
	if assetID == 0 {
		m.caption = "Create card"
		m.action = "CREATE"
		var t textinput.Model
		for i := range m.textInputs {
			t = textinput.New()
			t.Cursor.Style = cursorStyle
			switch i {
			case 0:
				t.Placeholder = "Sticker"
				t.CharLimit = 20
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle
				t.Focus()
			case 1:
				t.Placeholder = "Number"
				t.CharLimit = 16 //4 x 4 digits
			case 2:
				t.Placeholder = "Name"
				t.CharLimit = 20 //because

			case 3:
				t.Placeholder = "MM"
				t.CharLimit = 2
			case 4:
				t.Placeholder = "YY"
				t.CharLimit = 2
			case 5:
				t.Placeholder = "CVV"
				t.CharLimit = 3
			}
			m.textInputs[i] = t
		}
	} else {
		m.caption = "Edit card"
		m.action = "UPDATE"

		sticker := textinput.New()
		sticker.CharLimit = 20
		sticker.PromptStyle = focusedStyle
		sticker.TextStyle = focusedStyle
		sticker.Focus()

		number := textinput.New()
		number.CharLimit = 16

		name := textinput.New()
		name.CharLimit = 20

		month := textinput.New()
		month.CharLimit = 2

		year := textinput.New()
		year.CharLimit = 2

		cvv := textinput.New()
		cvv.CharLimit = 3

		//Get asset data
		m.asset, err = client.App.GetAsset(context.Background(), assetID)
		if err != nil {
			m.err = err
			number.Placeholder = "Get error"
			name.Placeholder = "Get error"
		} else {
			card := models.Card{}
			err = json.Unmarshal(m.asset.Body, &card)
			if err != nil {
				m.err = err
				number.Placeholder = "Convert error"
				name.Placeholder = "Convert error"
			} else {
				sticker.SetValue(m.asset.Sticker)
				number.SetValue(card.Number)
				name.SetValue(card.Name)
				month.SetValue(card.ExpMonth)
				year.SetValue(card.ExpYear)
				cvv.SetValue(card.CVV)
			}
		}
		m.textInputs[0] = sticker
		m.textInputs[1] = number
		m.textInputs[2] = name
		m.textInputs[3] = month
		m.textInputs[4] = year
		m.textInputs[5] = cvv
	}
	return m
}

func (m cardScreen) Init() tea.Cmd {
	return textinput.Blink
}
func (m cardScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	output := termenv.NewOutput(os.Stdout)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			screen := RootScreen()
			return screen, screen.Init()
		case tea.KeyCtrlC:
			output.ClearScreen()
			return m, tea.Quit

		}

		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.textInputs) {
				errMsg := validate(m.textInputs)
				if errMsg != "" {
					//Validation fails - return
					m.err = fmt.Errorf("Validation error: %v", errMsg)
					m.focusIndex = 0
					return m, nil
				} else {
					m.err = nil
				}
				//Parsing asset to json
				card := models.Card{
					Number:   m.textInputs[1].Value(),
					Name:     m.textInputs[2].Value(),
					ExpMonth: m.textInputs[3].Value(),
					ExpYear:  m.textInputs[4].Value(),
					CVV:      m.textInputs[5].Value(),
				}
				body, err := json.Marshal(card)
				if err != nil {
					m.err = fmt.Errorf("JSON error: %v", err)
					return m, nil
				}
				asset := &models.Asset{
					Type:    models.AssetTypeCard,
					Sticker: m.textInputs[0].Value(),
					Body:    body,
				}
				if m.action == "CREATE" {
					// Create new asset on server
					asset, err = client.App.CreateAsset(context.Background(), asset)
					if err != nil {
						m.focusIndex = 0
						m.err = fmt.Errorf("Asset create error: %v", err)
					} else {
						//Clear inputs
						for i := range m.textInputs {
							m.textInputs[i].SetValue("")
						}
						m.focusIndex = 0
						m.success = "Create successful"
						m.err = nil
					}
				} else if m.action == "UPDATE" {
					// Update asset on server
					asset.ID = m.asset.ID
					err = client.App.UpdateAsset(context.Background(), asset)
					if err != nil {
						m.focusIndex = 0
						m.err = fmt.Errorf("Asset update error: %v", err)
						m.success = "Update successful"
						m.err = nil
					}
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.textInputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.textInputs)
			}

			cmds := make([]tea.Cmd, len(m.textInputs))
			for i := 0; i <= len(m.textInputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.textInputs[i].Focus()
					m.textInputs[i].PromptStyle = focusedStyle
					m.textInputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.textInputs[i].Blur()
				m.textInputs[i].PromptStyle = noStyle
				m.textInputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *cardScreen) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))

	for i := range m.textInputs {
		m.textInputs[i], cmds[i] = m.textInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m cardScreen) View() string {
	var b strings.Builder
	b.WriteString(addKey.Render(m.caption))
	fmt.Fprintf(&b, "\n")

	// for i := range m.textInputs {
	// 	b.WriteString(viewStyle.Render(m.textInputs[i].View()))
	// 	if i < len(m.textInputs)-1 {
	// 		b.WriteRune('\n')
	// 	}
	// }
	//Stricker
	b.WriteString(viewStyle.Render(m.textInputs[0].View()))
	b.WriteRune('\n')

	//Number
	b.WriteString(viewStyle.Render(m.textInputs[1].View()))
	b.WriteRune('\n')
	//Name
	b.WriteString(viewStyle.Render(m.textInputs[2].View()))
	b.WriteRune('\n')
	//Month
	b.WriteString(viewStyle.Render(m.textInputs[3].View()))
	b.WriteRune('/')
	//Year
	b.WriteString(viewStyle.Render(m.textInputs[4].View()))
	b.WriteRune('\n')
	//CVV
	b.WriteString(viewStyle.Render(m.textInputs[5].View()))

	button := blurredButton.Render(m.action)
	if m.focusIndex == len(m.textInputs) {
		button = focusedButton.Render(m.action)
	}

	fmt.Fprintf(&b, "\n%s\n", button)

	if m.err != nil {
		fmt.Fprintf(&b, "\n%s\n", errorText.Render(m.err.Error()))
	}
	if m.success != "" {
		fmt.Fprintf(&b, "\n%s\n", successText.Render(m.success))
	}

	// Connection status 'widget'
	statusWidget(client.App.Workmode, &b)

	return b.String()
}

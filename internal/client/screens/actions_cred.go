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
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type credentialsScreen struct {
	focusIndex int
	textInputs []textinput.Model

	caption string
	action  string //for button caption

	err     error
	success string

	asset *models.Asset
}

func CreateCredentialsScreen(assetID int64) credentialsScreen {
	m := credentialsScreen{
		textInputs: make([]textinput.Model, 3),
	}
	var err error
	//Create new asset
	if assetID == 0 {
		m.caption = "Create login&pass"
		m.action = "CREATE"
		var t textinput.Model
		for i := range m.textInputs {
			t = textinput.New()
			t.Cursor.Style = cursorStyle
			t.CharLimit = 32

			switch i {
			case 0:
				t.Placeholder = "Sticker"
				t.Focus()
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle

			case 1:
				t.Placeholder = "Login"

			case 2:
				t.Placeholder = "Password"
				t.CharLimit = 64
				// t.EchoMode = textinput.EchoPassword
				// t.EchoCharacter = '*'

			}
			m.textInputs[i] = t
		}
	} else {
		m.caption = "Edit login&pass"
		m.action = "UPDATE"
		sticker := textinput.New()
		login := textinput.New()
		pass := textinput.New()
		//Get asset data
		m.asset, err = client.App.GetAsset(context.Background(), assetID)
		if err != nil {
			m.err = err
			sticker.Placeholder = "Get error"
			login.Placeholder = "Get error"
			pass.Placeholder = "Get error"
		} else {
			cred := models.Credentials{}
			err = json.Unmarshal(m.asset.Body, &cred)
			if err != nil {
				m.err = err
				sticker.Placeholder = "Convert error"
				login.Placeholder = "Convert error"
				pass.Placeholder = "Convert error"
			} else {
				sticker.SetValue(m.asset.Sticker)
				login.SetValue(cred.Login)
				pass.SetValue(cred.Password)
			}
		}
		m.textInputs[0] = sticker
		m.textInputs[1] = login
		m.textInputs[2] = pass
	}
	return m
}

func (m credentialsScreen) Init() tea.Cmd {
	return textinput.Blink
}
func (m credentialsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					m.err = fmt.Errorf("validation error: %v", errMsg)
					m.focusIndex = 0
					return m, nil
				} else {
					m.err = nil
				}
				//Parsing asset to json
				cred := models.Credentials{
					Login:    m.textInputs[1].Value(),
					Password: m.textInputs[2].Value(),
				}
				body, err := json.Marshal(cred)
				if err != nil {
					m.err = fmt.Errorf("json error: %v", err)
					return m, nil
				}
				asset := &models.Asset{
					Type:    models.AssetTypeCredentials,
					Sticker: m.textInputs[0].Value(),
					Body:    body,
				}
				if m.action == "CREATE" {
					// Create new asset on server
					_, err = client.App.CreateAsset(context.Background(), asset)
					if err != nil {
						m.focusIndex = 0
						m.err = fmt.Errorf("asset create error: %v", err)
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
						m.err = fmt.Errorf("asset update error: %v", err)
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

func (m *credentialsScreen) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))

	for i := range m.textInputs {
		m.textInputs[i], cmds[i] = m.textInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m credentialsScreen) View() string {
	var b strings.Builder
	b.WriteString(addKey.Render(m.caption))
	fmt.Fprintf(&b, "\n\n")

	for i := range m.textInputs {
		b.WriteString(viewStyle.Render(m.textInputs[i].View()))
		if i < len(m.textInputs)-1 {
			b.WriteRune('\n')
		}
	}

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

package screens

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"context"
	"fmt"
	"path/filepath"

	"os"
	"strings"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/Ssnakerss/mypreciouskeeper/internal/models"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type fileScreen struct {
	focusIndex int
	textInputs []textinput.Model
	cursorMode cursor.Mode

	caption string
	action  string //for button caption

	err     error
	success string

	asset *models.Asset

	buttonsQty int
}

func CreateFileScreen(assetID int64) fileScreen {
	m := fileScreen{}

	//Create new asset
	if assetID == 0 {
		m.caption = "Create file"
		m.action = "CREATE"

		//When crete file -  2 input fields - sticker and filename
		m.textInputs = make([]textinput.Model, 2)

		sticker := textinput.New()
		sticker.Cursor.Style = cursorStyle
		sticker.CharLimit = 32
		sticker.Placeholder = "Sticker"
		sticker.Focus()
		sticker.PromptStyle = focusedStyle
		sticker.TextStyle = focusedStyle
		m.textInputs[0] = sticker

		filename := textinput.New()
		filename.CharLimit = 100
		filename.Placeholder = "Select a file"
		m.textInputs[1] = filename

		//Number of action buttons
		//when create  - 1 button CREATE
		m.buttonsQty = 1
	} else {
		m.caption = "Edit file"
		m.action = "UPDATE"

		//When display/edit file - 3 input fields
		//Sticker, filename and path to save file to
		//And also 2 buttons - update and save file
		m.textInputs = make([]textinput.Model, 3)

		sticker := textinput.New()
		file := textinput.New()
		path := textinput.New()

		var err error

		//Get asset data
		m.asset, err = client.App.GetAsset(context.Background(), assetID)
		if err != nil {
			m.err = err
			sticker.Placeholder = "Get error"
			file.Placeholder = "Get error"
		} else {
			//Parse sticker into sticker and filename
			s := strings.Split(m.asset.Sticker, "|")
			if len(s) != 2 {
				m.err = fmt.Errorf("Asset sticker is not valid: %v", m.asset.Sticker)
				sticker.Placeholder = "Sticker is not valid"
				file.Placeholder = "Sticker is not valid"
			} else {
				sticker.SetValue(s[0])
				file.SetValue(s[1])
			}
		}
		m.textInputs[0] = sticker
		m.textInputs[1] = file
		m.textInputs[2] = path

		//Number of action buttons
		//when edit  - 2 buttons - update and save file
		m.buttonsQty = 2
	}

	return m
}

func (m fileScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (m fileScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				//Create new file asset
				errMsg := validate(m.textInputs)
				if errMsg != "" {
					//Validation fails - return
					m.err = fmt.Errorf("Validation error: %v", errMsg)
					m.focusIndex = 0
					return m, nil
				} else {
					m.err = nil
				}
				//File asset saving without json
				//Reading file into body directly
				//get filename with path to load
				filename := m.textInputs[1].Value()
				body, err := os.ReadFile(filename)
				if err != nil {
					m.focusIndex = 0
					m.err = fmt.Errorf("File read error: %v", err)
					return m, nil
				}
				//Get file name withot path
				filename = filepath.Base(filename)
				asset := &models.Asset{
					Type:    models.AssetTypeFile,
					Sticker: m.textInputs[0].Value() + "|" + filename,
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
			if s == "enter" && m.focusIndex == len(m.textInputs)+1 {
				//Save file to path
				path := m.textInputs[2].Value()
				filename := m.textInputs[1].Value()
				if m.asset != nil && m.asset.Type == models.AssetTypeFile {
					err := os.WriteFile(filepath.Join(path, filename), m.asset.Body, 0644)
					if err != nil {
						m.focusIndex = 0
						m.err = fmt.Errorf("File save error: %v", err)
					} else {
						m.focusIndex = 0
						m.success = "Save successful"
						m.err = nil
					}
				} else {
					m.focusIndex = 0
					m.err = fmt.Errorf("Asset is not file")
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++

			}

			if m.focusIndex > len(m.textInputs)+(m.buttonsQty-1) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.textInputs) + (m.buttonsQty - 1)
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

func (m *fileScreen) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))

	for i := range m.textInputs {
		m.textInputs[i], cmds[i] = m.textInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m fileScreen) View() string {
	var b strings.Builder
	b.WriteString(addKey.Render(m.caption))
	fmt.Fprintf(&b, "\n\n")

	for i := range m.textInputs {
		b.WriteString(viewStyle.Render(m.textInputs[i].View()))
		if i < len(m.textInputs)-1 {
			b.WriteRune('\n')
		}
	}

	var button string
	if m.focusIndex == len(m.textInputs) {
		button = focusedButton.Render(m.action)
	} else {
		button = blurredButton.Render(m.action)
	}
	fmt.Fprintf(&b, "\n%s\n", button)

	//If update - add button for file save to disk
	if m.action == "UPDATE" {
		if m.focusIndex == len(m.textInputs)+1 {
			button = focusedButton.Render("SAVE TO DISK")
		} else {
			button = blurredButton.Render("SAVE TO DISK")
		}
		fmt.Fprintf(&b, "\n%s\n", button)
	}

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

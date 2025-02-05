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
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type memoScreen struct {
	focusIndex int
	textInputs []textinput.Model
	cursorMode cursor.Mode

	textarea textarea.Model

	err     error
	success string
}

func CreateMemoScreen() memoScreen {
	ti := textarea.New()
	ti.Placeholder = "Your memo here  ..."

	m := memoScreen{
		textInputs: make([]textinput.Model, 1),
		textarea:   ti,
	}

	t := textinput.New()
	t.Cursor.Style = cursorStyle
	t.CharLimit = 32
	t.Placeholder = "Sticker"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	t.Validate = stickerValidator

	m.textInputs[0] = t

	return m
}

func (m memoScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (m memoScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	output := termenv.NewOutput(os.Stdout)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		//Quit the program
		case tea.KeyCtrlC:
			output.ClearScreen()
			return m, tea.Quit
			//Return to previous screen - Action menu
		case tea.KeyEsc, tea.KeyCtrlQ:
			screen_y := CreateActionsMenuScreen()
			return RootScreen().SwitchScreen(&screen_y)
		case tea.KeyTab:
			if m.textarea.Focused() {
				m.textarea.Blur()
				idx := len(m.textInputs) - 1
				m.focusIndex = idx
				m.textInputs[idx].Focus()
				m.textInputs[idx].PromptStyle = focusedStyle
				m.textInputs[idx].TextStyle = focusedStyle

			} else {
				m.textarea.Focus()

				m.textInputs[m.focusIndex].Blur()
				m.textInputs[m.focusIndex].PromptStyle = noStyle
				m.textInputs[m.focusIndex].TextStyle = noStyle

			}
		case tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			if m.textarea.Focused() {
				//For textarea enter is just a new row
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			} else {
				s := msg.String()
				//Handle Enter event on the button
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
					memo := models.Memo{
						Text: m.textarea.Value(),
					}
					body, err := json.Marshal(memo)
					if err != nil {
						m.err = fmt.Errorf("JSON error: %v", err)
						return m, nil
					}
					asset := &models.Asset{
						Type:    models.AssetTypeMemo,
						Sticker: m.textInputs[0].Value(),
						Body:    body,
					}

					// Create new asset on server
					asset, err = client.App.CreateAsset(context.Background(), asset)

					if err != nil {
						m.focusIndex = 0
						m.err = fmt.Errorf("Asset create error: %v", err)
					} else {
						//Clear inputs
						m.focusIndex = 0
						m.textarea.SetValue("")
						m.textInputs[0].SetValue("")
						m.success = "Create successful"
						m.err = nil
					}
				}

				//Cycle through text inputs
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
		default:
			// Processing text inputs
			if m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			} else {
				cmd = m.updateInputs(msg)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m *memoScreen) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))

	for i := range m.textInputs {
		m.textInputs[i], cmds[i] = m.textInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m memoScreen) View() string {
	var b strings.Builder
	b.WriteString(addKey.Render(fmt.Sprintf("New")))
	fmt.Fprintf(&b, "\n\n")

	for i := range m.textInputs {
		b.WriteString(viewStyle.Render(m.textInputs[i].View()))
		if i < len(m.textInputs)-1 {
			b.WriteRune('\n')
		}
	}

	fmt.Fprintf(&b,
		"\n\n%s\n\n",
		m.textarea.View(),
	)

	button := blurredButton.Render("[Create]")
	if m.focusIndex == len(m.textInputs) {
		button = focusedButton.Render("[Create]")
	}

	fmt.Fprintf(&b, "\n%s\n\n", button)

	if m.err != nil {
		fmt.Fprintf(&b, "\n%s\n", errorText.Render(m.err.Error()))
	}
	if m.success != "" {
		fmt.Fprintf(&b, "\n%s\n", successText.Render(m.success))
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", "(ctrl+c to quit)")

	return b.String()
}

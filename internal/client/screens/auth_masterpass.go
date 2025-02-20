package screens

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"

	"os"
	"strings"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type screenMasterPass struct {
	focusIndex int
	textInputs []textinput.Model

	err     error
	success string
	warning string
}

func ScreenMasterPass() screenMasterPass {
	m := screenMasterPass{
		textInputs: make([]textinput.Model, 1),
	}

	t := textinput.New()
	t.Cursor.Style = cursorStyle
	t.CharLimit = 32

	t.Placeholder = "your master pass"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	t.Validate = passwordValidator

	m.textInputs[0] = t

	return m
}

func (m screenMasterPass) Init() tea.Cmd {
	return textinput.Blink
}
func (m screenMasterPass) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.success = ""
				m.warning = ""
				errMsg := validate(m.textInputs)
				if errMsg != "" {
					//Validation fails - return
					m.err = fmt.Errorf("validation error: %v", errMsg)
					m.focusIndex = 0
					return m, nil
				} else {
					m.err = nil
				}
				//Try Login via gRPC
				client.App.SetMasterPass(m.textInputs[0].Value())
				m.success = fmt.Sprintf("Master pass is set to '%s'", string(client.App.GetMasterPass()))
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
func (m *screenMasterPass) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))

	for i := range m.textInputs {
		m.textInputs[i], cmds[i] = m.textInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Render screen view
func (m screenMasterPass) View() string {
	var b strings.Builder
	b.WriteString(addKey.Render("ENTER MASTER PASSWORD"))
	fmt.Fprintf(&b, "\n\n")

	for i := range m.textInputs {
		b.WriteString(viewStyle.Render(m.textInputs[i].View()))
		if i < len(m.textInputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := blurredButton.Render("[SAVE]")
	if m.focusIndex == len(m.textInputs) {
		button = focusedButton.Render("[SAVE]")
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", button)

	if m.err != nil {
		fmt.Fprintf(&b, "\n%s\n", errorText.Render(m.err.Error()))
	}
	if m.success != "" {
		fmt.Fprintf(&b, "\n%s\n", successText.Render(m.success))
	}
	if m.warning != "" {
		fmt.Fprintf(&b, "\n%s\n", warningText.Render(m.warning))
	}

	//Connection status 'widget'
	statusWidget(client.App.Workmode, &b)

	return b.String()
}

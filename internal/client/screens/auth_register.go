package screens

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"context"
	"fmt"

	"os"
	"strings"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

type screenRegister struct {
	focusIndex int
	textInputs []textinput.Model
	cursorMode cursor.Mode
	err        error
	error      string
	success    string
}

func ScreenRegister() screenRegister {
	m := screenRegister{
		textInputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.textInputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Email"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.Validate = emailValidator
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
			t.CharLimit = 64
			t.Validate = passwordValidator
		case 2:
			t.Placeholder = "Confirm password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '*'
			t.CharLimit = 64
		}
		m.textInputs[i] = t
	}

	return m
}

func (m screenRegister) Init() tea.Cmd {
	return textinput.Blink
}

func (m screenRegister) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			var err error
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

				//Before register check password match
				if m.textInputs[1].Value() != m.textInputs[2].Value() {
					m.focusIndex = 2
					m.err = fmt.Errorf("Password does not match")
					return m, nil
				}

				client.App.UserID, err = client.App.Register(context.Background(), m.textInputs[0].Value(), m.textInputs[1].Value())

				if err != nil {
					m.focusIndex = 1
					m.err = fmt.Errorf("Register error: %v", err)
				} else {
					m.success = "Register successful"
					m.err = nil
					// screen := RootScreen()
					// return RootScreen().SwitchScreen(&screen)
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
func (m *screenRegister) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))

	for i := range m.textInputs {
		m.textInputs[i], cmds[i] = m.textInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m screenRegister) View() string {
	var b strings.Builder
	b.WriteString(addKey.Render(fmt.Sprintf("Register")))

	fmt.Fprintf(&b, "\n\n")

	for i := range m.textInputs {
		b.WriteString(viewStyle.Render(m.textInputs[i].View()))
		if i < len(m.textInputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := blurredButton.Render("Register")
	if m.focusIndex == len(m.textInputs) {
		button = focusedButton.Render("Register")
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", button)

	if m.err != nil {
		fmt.Fprintf(&b, "\n%s\n", errorText.Render(m.err.Error()))
	}
	if m.success != "" {
		fmt.Fprintf(&b, "\n%s\n", successText.Render(m.success))
	}

	return b.String()
}

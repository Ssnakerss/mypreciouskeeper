package screens

import (
	"fmt"
	"io"
	"strings"

	client "github.com/Ssnakerss/mypreciouskeeper/internal/client/app"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//Styles

const defaultHeight = 20
const defaultWidth = 30

type (
	errMsg error
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	helpStyleInput      = blurredStyle

	focusedButton = focusedStyle.Padding(0, 2).Bold(true)
	blurredButton = blurredStyle.Padding(0, 2)

	errorText   = lipgloss.NewStyle().Padding(0, 2).Bold(true).Foreground(lipgloss.Color("#FF7575"))
	successText = lipgloss.NewStyle().Padding(0, 2).Bold(true).Foreground(lipgloss.Color("#00FF21"))
	warningText = lipgloss.NewStyle().Padding(0, 2).Bold(true).Foreground(lipgloss.Color("#FFFF00"))

	addKey    = lipgloss.NewStyle().MarginTop(1).MarginLeft(2).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#FFFFFF")).Padding(0, 5, 0, 5)
	viewStyle = lipgloss.NewStyle().Padding(0, 2)
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("#FFFFFF")).Padding(0, 5, 0, 5)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#54B575"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#A1FCC0"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	alias string
	title string
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("→ " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// Widget
func statusWidget(workMode string, b *strings.Builder) {
	if workMode == client.LOCAL {
		fmt.Fprintf(b, "\n%s", warningText.Render("Local mode"))
	} else {
		fmt.Fprintf(b, "\n%s", successText.Render("Remote mode"))
	}
}

// Validators
func emailValidator(s string) error {
	if len(s) < 5 || len(s) > 254 {
		return fmt.Errorf("email must be between 5 and 254 characters")
	}
	if !strings.Contains(s, "@") {
		return fmt.Errorf("email must contain @")
	}
	if !strings.Contains(s, ".") {
		return fmt.Errorf("email must contain .")
	}
	return nil
}
func passwordValidator(s string) error {
	if s == "" {
		return fmt.Errorf("password must not be empty")
	}
	return nil
}

func stickerValidator(s string) error {
	if s == "" {
		return fmt.Errorf("please input some sticker, please .....")
	}
	return nil
}

// Loop through inputs and call Validation func on each
func validate(inputs []textinput.Model) (errMsg string) {
	// Validation
	for i := range inputs {
		if inputs[i].Validate != nil {
			err := inputs[i].Validate(inputs[i].Value())
			if err != nil {
				errMsg += err.Error()
			}
		}
	}
	return errMsg
}

package screens

import (
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonkneen/pulsereader/internal/tui/components"
	"github.com/jasonkneen/pulsereader/internal/tui/styles"
)

// FileSelectedMsg is sent when a file is selected.
type FileSelectedMsg struct{ Path string }

// FilePickerCancelMsg is sent when the picker is cancelled.
type FilePickerCancelMsg struct{}

// FilePickerScreen wraps the bubbles filepicker.
type FilePickerScreen struct {
	picker filepicker.Model
	width  int
	height int
}

// NewFilePickerScreen creates a file picker for PDFs and EPUBs.
func NewFilePickerScreen(startDir string) FilePickerScreen {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".pdf", ".epub"}
	fp.CurrentDirectory = startDir
	fp.AutoHeight = false
	fp.Styles.Cursor = lipgloss.NewStyle().Foreground(styles.Purple)
	fp.Styles.Selected = lipgloss.NewStyle().Foreground(styles.Purple).Bold(true)

	return FilePickerScreen{
		picker: fp,
	}
}

func (s FilePickerScreen) Init() tea.Cmd {
	return s.picker.Init()
}

func (s FilePickerScreen) Update(msg tea.Msg) (FilePickerScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.picker.Height = msg.Height - 4

	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("esc", "q"))) {
			return s, func() tea.Msg { return FilePickerCancelMsg{} }
		}
	}

	var cmd tea.Cmd
	s.picker, cmd = s.picker.Update(msg)

	// Check if a file was selected
	if didSelect, path := s.picker.DidSelectFile(msg); didSelect {
		return s, func() tea.Msg { return FileSelectedMsg{Path: path} }
	}

	return s, cmd
}

func (s FilePickerScreen) View() string {
	title := styles.TitleStyle.
		Background(styles.Purple).
		Foreground(lipgloss.Color("#ffffff")).
		Padding(0, 2).
		Render("Open Document")

	help := components.HelpBar([]components.HelpBinding{
		{Key: "enter", Desc: "select"},
		{Key: "esc", Desc: "cancel"},
	})

	return lipgloss.JoinVertical(lipgloss.Left, title, "", s.picker.View(), help)
}

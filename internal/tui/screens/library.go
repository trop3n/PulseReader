package screens

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonkneen/pulsereader/internal/storage"
	"github.com/jasonkneen/pulsereader/internal/tui/components"
	"github.com/jasonkneen/pulsereader/internal/tui/styles"
)

// Messages for screen transitions.
type (
	OpenFilePickerMsg struct{}
	OpenDocumentMsg   struct{ Path string }
	RemoveDocumentMsg struct{ ID int64 }
	SwitchToFeedsMsg  struct{}
)

// libraryItem implements list.Item for the document list.
type libraryItem struct {
	doc storage.Document
}

func (i libraryItem) Title() string {
	icon := "📄"
	if i.doc.Format == "epub" {
		icon = "📖"
	}
	return fmt.Sprintf("%s %s", icon, i.doc.Title)
}

func (i libraryItem) Description() string {
	progress := ""
	if i.doc.TotalPages > 0 {
		pct := float64(i.doc.CurrentPage) / float64(i.doc.TotalPages) * 100
		if i.doc.CurrentPage > 0 {
			progress = fmt.Sprintf(" • %d/%d (%.0f%%)", i.doc.CurrentPage, i.doc.TotalPages, pct)
		} else {
			pageLabel := "pages"
			if i.doc.Format == "epub" {
				pageLabel = "chapters"
			}
			progress = fmt.Sprintf(" • %d %s", i.doc.TotalPages, pageLabel)
		}
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(i.doc.Format), progress)
}

func (i libraryItem) FilterValue() string { return i.doc.Title }

// LibraryScreen is the main document library view.
type LibraryScreen struct {
	list   list.Model
	db     *storage.DB
	width  int
	height int
	err    error
}

// NewLibraryScreen creates the library screen.
func NewLibraryScreen(db *storage.DB) LibraryScreen {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(styles.Purple).
		BorderLeftForeground(styles.Purple)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(styles.FgSecondary).
		BorderLeftForeground(styles.Purple)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "PulseReader"
	l.Styles.Title = styles.TitleStyle.
		Background(styles.Purple).
		Foreground(lipgloss.Color("#ffffff")).
		Padding(0, 2)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false) // We render our own help

	return LibraryScreen{
		list: l,
		db:   db,
	}
}

// LoadDocuments refreshes the document list from the database.
func (s *LibraryScreen) LoadDocuments() tea.Cmd {
	return func() tea.Msg {
		docs, err := s.db.ListDocuments()
		if err != nil {
			return errMsg{err}
		}
		return docsLoadedMsg{docs}
	}
}

type docsLoadedMsg struct{ docs []storage.Document }
type errMsg struct{ err error }

func (s LibraryScreen) Init() tea.Cmd {
	return s.LoadDocuments()
}

func (s LibraryScreen) Update(msg tea.Msg) (LibraryScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetSize(msg.Width, msg.Height-2) // Reserve space for help bar

	case docsLoadedMsg:
		items := make([]list.Item, len(msg.docs))
		for i, doc := range msg.docs {
			items[i] = libraryItem{doc: doc}
		}
		s.list.SetItems(items)
		return s, nil

	case errMsg:
		s.err = msg.err
		return s, nil

	case tea.KeyMsg:
		// Don't handle keys while filtering
		if s.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			return s, func() tea.Msg { return SwitchToFeedsMsg{} }

		case key.Matches(msg, key.NewBinding(key.WithKeys("a", "o"))):
			return s, func() tea.Msg { return OpenFilePickerMsg{} }

		case key.Matches(msg, key.NewBinding(key.WithKeys("b"))):
			return s, func() tea.Msg { return OpenFolderBrowserMsg{} }

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if item, ok := s.list.SelectedItem().(libraryItem); ok {
				return s, func() tea.Msg { return OpenDocumentMsg{Path: item.doc.Path} }
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("d", "delete"))):
			if item, ok := s.list.SelectedItem().(libraryItem); ok {
				return s, func() tea.Msg { return RemoveDocumentMsg{ID: item.doc.ID} }
			}
		}
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s LibraryScreen) View() string {
	help := components.HelpBar([]components.HelpBinding{
		{Key: "enter", Desc: "open"},
		{Key: "a/o", Desc: "add file"},
		{Key: "b", Desc: "browse folder"},
		{Key: "d", Desc: "remove"},
		{Key: "/", Desc: "filter"},
		{Key: "tab", Desc: "feeds"},
		{Key: "q", Desc: "quit"},
	})

	if s.err != nil {
		errStr := styles.ErrorStyle.Render("Error: " + s.err.Error())
		return lipgloss.JoinVertical(lipgloss.Left, s.list.View(), errStr, help)
	}

	return lipgloss.JoinVertical(lipgloss.Left, s.list.View(), help)
}

// AddDocumentFromPath adds a document to the library from a file path.
func (s *LibraryScreen) AddDocumentFromPath(path string) tea.Cmd {
	return func() tea.Msg {
		// Determine format
		ext := strings.ToLower(filepath.Ext(path))
		format := ""
		switch ext {
		case ".pdf":
			format = "pdf"
		case ".epub":
			format = "epub"
		default:
			return errMsg{fmt.Errorf("unsupported format: %s", ext)}
		}

		title := strings.TrimSuffix(filepath.Base(path), ext)

		doc := &storage.Document{
			Title:  title,
			Path:   path,
			Format: format,
		}
		if err := s.db.AddDocument(doc); err != nil {
			return errMsg{err}
		}
		// Reload the list
		docs, err := s.db.ListDocuments()
		if err != nil {
			return errMsg{err}
		}
		return docsLoadedMsg{docs}
	}
}

package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonkneen/pulsereader/internal/document"
	"github.com/jasonkneen/pulsereader/internal/storage"
	"github.com/jasonkneen/pulsereader/internal/tui/components"
	"github.com/jasonkneen/pulsereader/internal/tui/styles"
)

// CloseReaderMsg signals returning to the library.
type CloseReaderMsg struct{}

// ReaderScreen displays document content in a scrollable viewport.
type ReaderScreen struct {
	doc      document.Document
	dbDoc    *storage.Document
	db       *storage.DB
	viewport viewport.Model
	page     int // 1-indexed current page/chapter
	content  string
	width    int
	height   int
	ready    bool
	err      error
	renderer *glamour.TermRenderer
}

// NewReaderScreen creates a reader for the given document.
func NewReaderScreen(doc document.Document, dbDoc *storage.Document, db *storage.DB, width, height int) ReaderScreen {
	vp := viewport.New(width, height-3) // Reserve for title + status + help
	vp.Style = styles.ReaderContentStyle

	startPage := 1
	if dbDoc != nil && dbDoc.CurrentPage > 0 {
		startPage = dbDoc.CurrentPage
	}

	// Create glamour renderer for markdown content
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-10), // Account for padding
	)

	s := ReaderScreen{
		doc:      doc,
		dbDoc:    dbDoc,
		db:       db,
		viewport: vp,
		page:     startPage,
		width:    width,
		height:   height,
		ready:    true,
		renderer: renderer,
	}
	s.loadPage()
	return s
}

func (s *ReaderScreen) loadPage() {
	content, err := s.doc.PageContent(s.page)
	if err != nil {
		s.err = err
		s.content = ""
		return
	}

	// For EPUBs, the content is already markdown — render it with glamour
	meta := s.doc.Metadata()
	if meta.Format == "epub" && s.renderer != nil {
		rendered, err := s.renderer.Render(content)
		if err == nil {
			content = rendered
		}
	}

	s.content = content
	s.viewport.SetContent(content)
	s.viewport.GotoTop()
	s.err = nil
}

func (s *ReaderScreen) saveProgress() {
	if s.dbDoc != nil && s.db != nil {
		pct := s.viewport.ScrollPercent()
		_ = s.db.SaveProgress(s.dbDoc.ID, s.page, pct)
	}
}

func (s ReaderScreen) Init() tea.Cmd {
	return nil
}

func (s ReaderScreen) Update(msg tea.Msg) (ReaderScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.viewport.Width = msg.Width
		s.viewport.Height = msg.Height - 3

		// Recreate renderer with new width
		s.renderer, _ = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(msg.Width-10),
		)
		s.loadPage()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc"))):
			s.saveProgress()
			return s, func() tea.Msg { return CloseReaderMsg{} }

		// Next page/chapter
		case key.Matches(msg, key.NewBinding(key.WithKeys("right", "l", "]"))):
			if s.page < s.doc.TotalPages() {
				s.page++
				s.loadPage()
				s.saveProgress()
			}
			return s, nil

		// Previous page/chapter
		case key.Matches(msg, key.NewBinding(key.WithKeys("left", "h", "["))):
			if s.page > 1 {
				s.page--
				s.loadPage()
				s.saveProgress()
			}
			return s, nil

		// Jump to first page
		case key.Matches(msg, key.NewBinding(key.WithKeys("home"))):
			s.page = 1
			s.loadPage()
			return s, nil

		// Jump to last page
		case key.Matches(msg, key.NewBinding(key.WithKeys("end"))):
			s.page = s.doc.TotalPages()
			s.loadPage()
			return s, nil
		}
	}

	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s ReaderScreen) View() string {
	meta := s.doc.Metadata()

	// Title bar
	pageLabel := "Page"
	if meta.Format == "epub" {
		pageLabel = "Chapter"
	}
	titleText := styles.AccentStyle.Render(meta.Title)
	if meta.Author != "" {
		titleText += styles.SubtitleStyle.Render(" by " + meta.Author)
	}
	titleBar := lipgloss.NewStyle().
		Width(s.width).
		Background(styles.BgPanel).
		Padding(0, 1).
		Render(titleText)

	// Error display
	if s.err != nil {
		errContent := styles.ErrorStyle.Render(fmt.Sprintf("Error loading content: %v", s.err))
		return lipgloss.JoinVertical(lipgloss.Left, titleBar, errContent)
	}

	// Status bar
	scrollPct := s.viewport.ScrollPercent() * 100
	status := components.StatusBar(
		s.width,
		fmt.Sprintf("%s %d/%d", pageLabel, s.page, s.doc.TotalPages()),
		fmt.Sprintf("%.0f%%", scrollPct),
		meta.Format,
	)

	// Help
	helpBindings := []components.HelpBinding{
		{Key: "j/k", Desc: "scroll"},
		{Key: "←/→", Desc: pageLabel},
		{Key: "q", Desc: "back"},
	}
	help := components.HelpBar(helpBindings)

	// Build content - if empty, show a message
	viewContent := s.viewport.View()
	if strings.TrimSpace(s.content) == "" {
		viewContent = lipgloss.Place(
			s.width, s.height-3,
			lipgloss.Center, lipgloss.Center,
			styles.SubtitleStyle.Render("(no text content on this page)"),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, titleBar, viewContent, status, help)
}

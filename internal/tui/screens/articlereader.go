package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonkneen/pulsereader/internal/feed"
	"github.com/jasonkneen/pulsereader/internal/storage"
	"github.com/jasonkneen/pulsereader/internal/tui/components"
	"github.com/jasonkneen/pulsereader/internal/tui/styles"
)

// CloseArticleReaderMsg signals returning to the article list.
type CloseArticleReaderMsg struct{}

// ArticleReaderScreen displays an article in a scrollable viewport.
type ArticleReaderScreen struct {
	article  *storage.Article
	db       *storage.DB
	viewport viewport.Model
	content  string
	width    int
	height   int
	renderer *glamour.TermRenderer
}

// NewArticleReaderScreen creates a reader for the given article.
func NewArticleReaderScreen(article *storage.Article, db *storage.DB, width, height int) ArticleReaderScreen {
	vp := viewport.New(width, height-3)
	vp.Style = styles.ReaderContentStyle

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-10),
	)

	// Convert article content to markdown
	markdown := feed.RenderArticle(article)

	// Render with glamour
	content := markdown
	if renderer != nil {
		rendered, err := renderer.Render(markdown)
		if err == nil {
			content = rendered
		}
	}

	vp.SetContent(content)

	// Mark as read
	_ = db.MarkArticleRead(article.ID)

	return ArticleReaderScreen{
		article:  article,
		db:       db,
		viewport: vp,
		content:  content,
		width:    width,
		height:   height,
		renderer: renderer,
	}
}

func (s ArticleReaderScreen) Init() tea.Cmd {
	return nil
}

func (s ArticleReaderScreen) Update(msg tea.Msg) (ArticleReaderScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.viewport.Width = msg.Width
		s.viewport.Height = msg.Height - 3

		s.renderer, _ = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(msg.Width-10),
		)
		markdown := feed.RenderArticle(s.article)
		if s.renderer != nil {
			rendered, err := s.renderer.Render(markdown)
			if err == nil {
				s.content = rendered
				s.viewport.SetContent(s.content)
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc"))):
			return s, func() tea.Msg { return CloseArticleReaderMsg{} }
		}
	}

	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return s, cmd
}

func (s ArticleReaderScreen) View() string {
	// Title bar
	titleText := styles.AccentStyle.Render(s.article.Title)
	if s.article.FeedTitle != "" {
		titleText += styles.SubtitleStyle.Render(" — " + s.article.FeedTitle)
	}
	titleBar := lipgloss.NewStyle().
		Width(s.width).
		Background(styles.BgPanel).
		Padding(0, 1).
		Render(titleText)

	// Status bar
	scrollPct := s.viewport.ScrollPercent() * 100
	status := components.StatusBar(
		s.width,
		s.article.FeedTitle,
		fmt.Sprintf("%.0f%%", scrollPct),
		"",
	)

	help := components.HelpBar([]components.HelpBinding{
		{Key: "j/k", Desc: "scroll"},
		{Key: "q", Desc: "back"},
	})

	viewContent := s.viewport.View()
	if strings.TrimSpace(s.content) == "" {
		viewContent = lipgloss.Place(
			s.width, s.height-3,
			lipgloss.Center, lipgloss.Center,
			styles.SubtitleStyle.Render("(no content available)"),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, titleBar, viewContent, status, help)
}

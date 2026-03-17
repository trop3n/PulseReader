package screens

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonkneen/pulsereader/internal/storage"
	"github.com/jasonkneen/pulsereader/internal/tui/components"
	"github.com/jasonkneen/pulsereader/internal/tui/styles"
)

// Article screen messages.
type (
	CloseArticleListMsg struct{}
	OpenArticleMsg      struct{ ArticleID int64 }
)

// articleItem implements list.Item.
type articleItem struct {
	article storage.Article
}

func (i articleItem) Title() string {
	prefix := "●"
	if i.article.Read {
		prefix = "○"
	}
	return fmt.Sprintf("%s %s", prefix, i.article.Title)
}

func (i articleItem) Description() string {
	parts := []string{}
	if i.article.FeedTitle != "" {
		parts = append(parts, i.article.FeedTitle)
	}
	if i.article.Author != "" {
		parts = append(parts, i.article.Author)
	}
	if i.article.PublishedAt != nil {
		parts = append(parts, timeAgo(*i.article.PublishedAt))
	}
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for _, p := range parts[1:] {
		result += " · " + p
	}
	return result
}

func (i articleItem) FilterValue() string { return i.article.Title }

type articlesLoadedMsg struct{ articles []storage.Article }

// ArticleListScreen shows articles for a feed or all feeds.
type ArticleListScreen struct {
	list       list.Model
	db         *storage.DB
	feedID     int64  // 0 = all feeds
	feedTitle  string
	unreadOnly bool
	width      int
	height     int
	err        error
}

// NewArticleListScreen creates an article list screen.
func NewArticleListScreen(db *storage.DB, feedID int64, feedTitle string) ArticleListScreen {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(styles.Purple).
		BorderLeftForeground(styles.Purple)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(styles.FgSecondary).
		BorderLeftForeground(styles.Purple)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = feedTitle
	l.Styles.Title = styles.TitleStyle.
		Background(styles.PurpleDim).
		Foreground(lipgloss.Color("#ffffff")).
		Padding(0, 2)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return ArticleListScreen{
		list:      l,
		db:        db,
		feedID:    feedID,
		feedTitle: feedTitle,
	}
}

// LoadArticles fetches articles from the database.
func (s *ArticleListScreen) LoadArticles() tea.Cmd {
	feedID := s.feedID
	unreadOnly := s.unreadOnly
	db := s.db
	return func() tea.Msg {
		articles, err := db.ListArticles(feedID, unreadOnly, 200, 0)
		if err != nil {
			return feedErrMsg{err}
		}
		return articlesLoadedMsg{articles}
	}
}

func (s ArticleListScreen) Init() tea.Cmd {
	return s.LoadArticles()
}

func (s ArticleListScreen) Update(msg tea.Msg) (ArticleListScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetSize(msg.Width, msg.Height-2)

	case articlesLoadedMsg:
		items := make([]list.Item, len(msg.articles))
		for i, a := range msg.articles {
			items[i] = articleItem{article: a}
		}
		s.list.SetItems(items)
		s.err = nil
		return s, nil

	case feedErrMsg:
		s.err = msg.err
		return s, nil

	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc"))):
			return s, func() tea.Msg { return CloseArticleListMsg{} }

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if item, ok := s.list.SelectedItem().(articleItem); ok {
				id := item.article.ID
				return s, func() tea.Msg { return OpenArticleMsg{ArticleID: id} }
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("u"))):
			// Toggle unread only
			s.unreadOnly = !s.unreadOnly
			title := s.feedTitle
			if s.unreadOnly {
				title += " (unread)"
			}
			s.list.Title = title
			return s, s.LoadArticles()

		case key.Matches(msg, key.NewBinding(key.WithKeys("m"))):
			// Mark all read
			_ = s.db.MarkAllRead(s.feedID)
			return s, s.LoadArticles()
		}
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s ArticleListScreen) View() string {
	help := components.HelpBar([]components.HelpBinding{
		{Key: "enter", Desc: "read"},
		{Key: "u", Desc: "toggle unread"},
		{Key: "m", Desc: "mark all read"},
		{Key: "/", Desc: "filter"},
		{Key: "q", Desc: "back"},
	})

	if s.err != nil {
		errStr := styles.ErrorStyle.Render("Error: " + s.err.Error())
		return lipgloss.JoinVertical(lipgloss.Left, s.list.View(), errStr, help)
	}

	return lipgloss.JoinVertical(lipgloss.Left, s.list.View(), help)
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%dh ago", h)
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%dd ago", days)
	default:
		return t.Format("Jan 2")
	}
}

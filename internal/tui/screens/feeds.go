package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonkneen/pulsereader/internal/feed"
	"github.com/jasonkneen/pulsereader/internal/storage"
	"github.com/jasonkneen/pulsereader/internal/tui/components"
	"github.com/jasonkneen/pulsereader/internal/tui/styles"
)

// Feed screen messages.
type (
	SwitchToLibraryMsg  struct{}
	OpenFeedMsg         struct{ FeedID int64 }
	OpenAllArticlesMsg  struct{}
	RemoveFeedMsg       struct{ ID int64 }
)

// feedItem implements list.Item.
type feedItem struct {
	feed storage.Feed
}

func (i feedItem) Title() string {
	unread := ""
	if i.feed.UnreadCount > 0 {
		unread = fmt.Sprintf(" (%d)", i.feed.UnreadCount)
	}
	return fmt.Sprintf("📡 %s%s", i.feed.Title, unread)
}

func (i feedItem) Description() string {
	if i.feed.Description != "" {
		desc := i.feed.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		return desc
	}
	return i.feed.URL
}

func (i feedItem) FilterValue() string { return i.feed.Title }

// Internal messages.
type (
	feedsLoadedMsg   struct{ feeds []storage.Feed }
	feedErrMsg       struct{ err error }
	feedAddedMsg     struct{ result feed.FetchResult }
	feedsRefreshedMsg struct{ results []feed.FetchResult }
)

// FeedsScreen shows the list of subscribed feeds.
type FeedsScreen struct {
	list      list.Model
	db        *storage.DB
	width     int
	height    int
	err       error
	adding    bool
	urlInput  textinput.Model
	fetching  bool
}

// NewFeedsScreen creates the feeds screen.
func NewFeedsScreen(db *storage.DB) FeedsScreen {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(styles.Purple).
		BorderLeftForeground(styles.Purple)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(styles.FgSecondary).
		BorderLeftForeground(styles.Purple)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Feeds"
	l.Styles.Title = styles.TitleStyle.
		Background(styles.Purple).
		Foreground(lipgloss.Color("#ffffff")).
		Padding(0, 2)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	ti := textinput.New()
	ti.Placeholder = "https://example.com/feed.xml"
	ti.CharLimit = 500
	ti.Width = 60

	return FeedsScreen{
		list:     l,
		db:       db,
		urlInput: ti,
	}
}

// LoadFeeds refreshes the feed list.
func (s *FeedsScreen) LoadFeeds() tea.Cmd {
	return func() tea.Msg {
		feeds, err := s.db.ListFeeds()
		if err != nil {
			return feedErrMsg{err}
		}
		return feedsLoadedMsg{feeds}
	}
}

func (s FeedsScreen) Init() tea.Cmd {
	return s.LoadFeeds()
}

func (s FeedsScreen) Update(msg tea.Msg) (FeedsScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetSize(msg.Width, msg.Height-2)

	case feedsLoadedMsg:
		// Build items: "All Articles" entry first, then individual feeds
		items := make([]list.Item, 0, len(msg.feeds)+1)
		totalUnread := 0
		for _, f := range msg.feeds {
			totalUnread += f.UnreadCount
		}
		items = append(items, allArticlesItem{unread: totalUnread})
		for _, f := range msg.feeds {
			items = append(items, feedItem{feed: f})
		}
		s.list.SetItems(items)
		s.err = nil
		return s, nil

	case feedErrMsg:
		s.err = msg.err
		return s, nil

	case feedAddedMsg:
		s.adding = false
		s.fetching = false
		if msg.result.Err != nil {
			s.err = msg.result.Err
			return s, nil
		}
		return s, s.LoadFeeds()

	case feedsRefreshedMsg:
		s.fetching = false
		// Check for errors
		for _, r := range msg.results {
			if r.Err != nil {
				s.err = r.Err
				break
			}
		}
		return s, s.LoadFeeds()

	case tea.KeyMsg:
		// Handle URL input mode
		if s.adding {
			switch msg.String() {
			case "enter":
				url := s.urlInput.Value()
				if url != "" {
					s.fetching = true
					s.urlInput.SetValue("")
					return s, s.addFeed(url)
				}
				s.adding = false
				return s, nil
			case "esc":
				s.adding = false
				s.urlInput.SetValue("")
				return s, nil
			}
			var cmd tea.Cmd
			s.urlInput, cmd = s.urlInput.Update(msg)
			return s, cmd
		}

		// Don't handle keys while filtering
		if s.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			return s, func() tea.Msg { return SwitchToLibraryMsg{} }

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if item, ok := s.list.SelectedItem().(feedItem); ok {
				id := item.feed.ID
				return s, func() tea.Msg { return OpenFeedMsg{FeedID: id} }
			}
			if _, ok := s.list.SelectedItem().(allArticlesItem); ok {
				return s, func() tea.Msg { return OpenAllArticlesMsg{} }
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("a"))):
			s.adding = true
			s.err = nil
			s.urlInput.Focus()
			return s, textinput.Blink

		case key.Matches(msg, key.NewBinding(key.WithKeys("r", "R"))):
			if !s.fetching {
				s.fetching = true
				return s, s.refreshAll()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("d", "delete"))):
			if item, ok := s.list.SelectedItem().(feedItem); ok {
				id := item.feed.ID
				return s, func() tea.Msg { return RemoveFeedMsg{ID: id} }
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("m"))):
			// Mark all read in selected feed
			if item, ok := s.list.SelectedItem().(feedItem); ok {
				_ = s.db.MarkAllRead(item.feed.ID)
				return s, s.LoadFeeds()
			}
			if _, ok := s.list.SelectedItem().(allArticlesItem); ok {
				_ = s.db.MarkAllRead(0)
				return s, s.LoadFeeds()
			}
		}
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s FeedsScreen) View() string {
	var content string

	if s.adding {
		title := styles.TitleStyle.
			Background(styles.Purple).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 2).
			Render("Add Feed")
		prompt := lipgloss.NewStyle().Padding(1, 2).Render(
			styles.SubtitleStyle.Render("Feed URL: ") + s.urlInput.View(),
		)
		help := components.HelpBar([]components.HelpBinding{
			{Key: "enter", Desc: "add"},
			{Key: "esc", Desc: "cancel"},
		})
		return lipgloss.JoinVertical(lipgloss.Left, title, prompt, help)
	}

	content = s.list.View()

	statusParts := []string{}
	if s.fetching {
		statusParts = append(statusParts, styles.AccentStyle.Render("refreshing..."))
	}
	if s.err != nil {
		statusParts = append(statusParts, styles.ErrorStyle.Render("Error: "+s.err.Error()))
	}

	help := components.HelpBar([]components.HelpBinding{
		{Key: "enter", Desc: "open"},
		{Key: "a", Desc: "add feed"},
		{Key: "r", Desc: "refresh"},
		{Key: "m", Desc: "mark read"},
		{Key: "d", Desc: "remove"},
		{Key: "tab", Desc: "documents"},
		{Key: "q", Desc: "quit"},
	})

	parts := []string{content}
	for _, sp := range statusParts {
		parts = append(parts, sp)
	}
	parts = append(parts, help)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (s *FeedsScreen) addFeed(url string) tea.Cmd {
	db := s.db
	return func() tea.Msg {
		result := feed.AddAndFetchFeed(db, url)
		return feedAddedMsg{result}
	}
}

func (s *FeedsScreen) refreshAll() tea.Cmd {
	db := s.db
	return func() tea.Msg {
		results := feed.FetchAllFeeds(db)
		return feedsRefreshedMsg{results}
	}
}

// allArticlesItem is the "All Articles" entry at the top of the feed list.
type allArticlesItem struct {
	unread int
}

func (i allArticlesItem) Title() string {
	unread := ""
	if i.unread > 0 {
		unread = fmt.Sprintf(" (%d)", i.unread)
	}
	return fmt.Sprintf("📰 All Articles%s", unread)
}
func (i allArticlesItem) Description() string { return "View articles from all feeds" }
func (i allArticlesItem) FilterValue() string { return "All Articles" }

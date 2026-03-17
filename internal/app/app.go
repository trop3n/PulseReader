package app

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jasonkneen/pulsereader/internal/document"
	"github.com/jasonkneen/pulsereader/internal/storage"
	"github.com/jasonkneen/pulsereader/internal/tui/screens"
)

type screen int

const (
	libraryScreen screen = iota
	readerScreen
	filePickerScreen
	folderBrowserScreen
	feedsScreen
	articleListScreen
	articleReaderScreen
)

// Model is the top-level application model.
type Model struct {
	db            *storage.DB
	active        screen
	library       screens.LibraryScreen
	reader        screens.ReaderScreen
	filePicker    screens.FilePickerScreen
	folderBrowser screens.FolderBrowserScreen
	feeds         screens.FeedsScreen
	articleList   screens.ArticleListScreen
	articleReader screens.ArticleReaderScreen
	width         int
	height        int
	openDoc       document.Document
	lastBrowseDir string
	// Track where to return from article reader
	returnToFeeds bool
}

// New creates the application model.
func New(db *storage.DB) Model {
	return Model{
		db:      db,
		active:  libraryScreen,
		library: screens.NewLibraryScreen(db),
		feeds:   screens.NewFeedsScreen(db),
	}
}

func (m Model) Init() tea.Cmd {
	return m.library.LoadDocuments()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// Quit from library or feeds with 'q' (only when not filtering)
		if msg.String() == "q" && (m.active == libraryScreen || m.active == feedsScreen) {
			// Let the screen handle it first (it may be filtering)
		}

	// === Tab switching between Library and Feeds ===
	case screens.SwitchToFeedsMsg:
		m.active = feedsScreen
		return m, m.feeds.LoadFeeds()

	case screens.SwitchToLibraryMsg:
		m.active = libraryScreen
		return m, m.library.LoadDocuments()

	// === File picker / folder browser ===
	case screens.OpenFilePickerMsg:
		homeDir, _ := os.UserHomeDir()
		m.filePicker = screens.NewFilePickerScreen(homeDir)
		m.active = filePickerScreen
		return m, m.filePicker.Init()

	case screens.OpenFolderBrowserMsg:
		startDir := msg.Dir
		if startDir == "" {
			startDir = m.lastBrowseDir
		}
		if startDir == "" {
			startDir = m.db.GetSetting("last_browse_dir")
		}
		if startDir == "" {
			startDir, _ = os.UserHomeDir()
		}
		m.folderBrowser = screens.NewFolderBrowserScreen(startDir)
		m.active = folderBrowserScreen
		return m, m.folderBrowser.LoadDir()

	case screens.CloseFolderBrowserMsg:
		m.lastBrowseDir = m.folderBrowser.Dir()
		_ = m.db.SetSetting("last_browse_dir", m.lastBrowseDir)
		m.active = libraryScreen
		return m, nil

	case screens.FileSelectedMsg:
		if m.active == folderBrowserScreen {
			m.lastBrowseDir = m.folderBrowser.Dir()
			_ = m.db.SetSetting("last_browse_dir", m.lastBrowseDir)
		}
		m.active = libraryScreen
		return m, m.library.AddDocumentFromPath(msg.Path)

	case screens.FilePickerCancelMsg:
		m.active = libraryScreen
		return m, nil

	// === Document reader ===
	case screens.OpenDocumentMsg:
		return m, m.openDocument(msg.Path)

	case documentOpenedMsg:
		m.openDoc = msg.doc
		m.reader = screens.NewReaderScreen(msg.doc, msg.dbDoc, m.db, m.width, m.height)
		m.active = readerScreen
		meta := msg.doc.Metadata()
		if msg.dbDoc != nil && msg.dbDoc.TotalPages != meta.TotalPages {
			msg.dbDoc.TotalPages = meta.TotalPages
			_ = m.db.AddDocument(msg.dbDoc)
		}
		_ = m.db.TouchDocument(msg.dbDoc.ID)
		return m, nil

	case documentOpenErrorMsg:
		return m, nil

	case screens.RemoveDocumentMsg:
		_ = m.db.RemoveDocument(msg.ID)
		return m, m.library.LoadDocuments()

	case screens.CloseReaderMsg:
		if m.openDoc != nil {
			m.openDoc.Close()
			m.openDoc = nil
		}
		m.active = libraryScreen
		return m, m.library.LoadDocuments()

	// === Feed management ===
	case screens.OpenFeedMsg:
		feed, err := m.db.GetFeed(msg.FeedID)
		if err != nil {
			return m, nil
		}
		m.articleList = screens.NewArticleListScreen(m.db, feed.ID, feed.Title)
		m.active = articleListScreen
		return m, m.articleList.LoadArticles()

	case screens.OpenAllArticlesMsg:
		m.articleList = screens.NewArticleListScreen(m.db, 0, "All Articles")
		m.active = articleListScreen
		return m, m.articleList.LoadArticles()

	case screens.RemoveFeedMsg:
		_ = m.db.RemoveFeed(msg.ID)
		return m, m.feeds.LoadFeeds()

	case screens.CloseArticleListMsg:
		m.active = feedsScreen
		return m, m.feeds.LoadFeeds()

	// === Article reader ===
	case screens.OpenArticleMsg:
		article, err := m.db.GetArticle(msg.ArticleID)
		if err != nil {
			return m, nil
		}
		m.articleReader = screens.NewArticleReaderScreen(article, m.db, m.width, m.height)
		m.active = articleReaderScreen
		return m, nil

	case screens.CloseArticleReaderMsg:
		m.active = articleListScreen
		return m, m.articleList.LoadArticles()
	}

	// Route to active screen
	var cmd tea.Cmd
	switch m.active {
	case libraryScreen:
		m.library, cmd = m.library.Update(msg)
	case readerScreen:
		m.reader, cmd = m.reader.Update(msg)
	case filePickerScreen:
		m.filePicker, cmd = m.filePicker.Update(msg)
	case folderBrowserScreen:
		m.folderBrowser, cmd = m.folderBrowser.Update(msg)
	case feedsScreen:
		m.feeds, cmd = m.feeds.Update(msg)
	case articleListScreen:
		m.articleList, cmd = m.articleList.Update(msg)
	case articleReaderScreen:
		m.articleReader, cmd = m.articleReader.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.active {
	case readerScreen:
		return m.reader.View()
	case filePickerScreen:
		return m.filePicker.View()
	case folderBrowserScreen:
		return m.folderBrowser.View()
	case feedsScreen:
		return m.feeds.View()
	case articleListScreen:
		return m.articleList.View()
	case articleReaderScreen:
		return m.articleReader.View()
	default:
		return m.library.View()
	}
}

type documentOpenedMsg struct {
	doc   document.Document
	dbDoc *storage.Document
}
type documentOpenErrorMsg struct{ err error }

func (m *Model) openDocument(path string) tea.Cmd {
	return func() tea.Msg {
		doc, err := document.Open(path)
		if err != nil {
			return documentOpenErrorMsg{err}
		}

		dbDoc, _ := m.db.GetDocumentByPath(path)
		if dbDoc == nil {
			meta := doc.Metadata()
			dbDoc = &storage.Document{
				Title:      meta.Title,
				Path:       path,
				Format:     meta.Format,
				Author:     meta.Author,
				TotalPages: meta.TotalPages,
				FileSize:   meta.FileSize,
			}
			_ = m.db.AddDocument(dbDoc)
			dbDoc, _ = m.db.GetDocumentByPath(path)
		}

		return documentOpenedMsg{doc: doc, dbDoc: dbDoc}
	}
}

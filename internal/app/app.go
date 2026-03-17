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
)

// Model is the top-level application model.
type Model struct {
	db            *storage.DB
	active        screen
	library       screens.LibraryScreen
	reader        screens.ReaderScreen
	filePicker    screens.FilePickerScreen
	folderBrowser screens.FolderBrowserScreen
	width         int
	height        int
	openDoc       document.Document
	lastBrowseDir string
}

// New creates the application model.
func New(db *storage.DB) Model {
	return Model{
		db:      db,
		active:  libraryScreen,
		library: screens.NewLibraryScreen(db),
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
		// Global quit from library
		if m.active == libraryScreen && msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.active == libraryScreen && msg.String() == "q" {
			// Only quit if not filtering
			// The library list handles 'q' during filtering
		}

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
		// Remember the directory we were browsing
		if m.active == folderBrowserScreen {
			m.lastBrowseDir = m.folderBrowser.Dir()
			_ = m.db.SetSetting("last_browse_dir", m.lastBrowseDir)
		}
		m.active = libraryScreen
		return m, m.library.AddDocumentFromPath(msg.Path)

	case screens.FilePickerCancelMsg:
		m.active = libraryScreen
		return m, nil

	case screens.OpenDocumentMsg:
		return m, m.openDocument(msg.Path)

	case documentOpenedMsg:
		m.openDoc = msg.doc
		m.reader = screens.NewReaderScreen(msg.doc, msg.dbDoc, m.db, m.width, m.height)
		m.active = readerScreen
		// Update document metadata in DB
		meta := msg.doc.Metadata()
		if msg.dbDoc != nil && msg.dbDoc.TotalPages != meta.TotalPages {
			msg.dbDoc.TotalPages = meta.TotalPages
			_ = m.db.AddDocument(msg.dbDoc)
		}
		_ = m.db.TouchDocument(msg.dbDoc.ID)
		return m, nil

	case documentOpenErrorMsg:
		// Stay on library, error will be shown
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
			// Shouldn't happen if document was added, but handle gracefully
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

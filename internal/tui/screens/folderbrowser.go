package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonkneen/pulsereader/internal/tui/components"
	"github.com/jasonkneen/pulsereader/internal/tui/styles"
)

// Messages
type (
	OpenFolderBrowserMsg struct{ Dir string }
	CloseFolderBrowserMsg struct{}
)

// folderEntry represents a directory or document in the browser.
type folderEntry struct {
	name  string
	path  string
	isDir bool
	size  int64
	ext   string // ".pdf", ".epub", or "" for dirs
}

func (e folderEntry) Title() string {
	if e.isDir {
		return "📁 " + e.name
	}
	icon := "📄"
	if e.ext == ".epub" {
		icon = "📖"
	}
	return icon + " " + e.name
}

func (e folderEntry) Description() string {
	if e.isDir {
		return "directory"
	}
	size := formatSize(e.size)
	return fmt.Sprintf("%s  %s", strings.ToUpper(strings.TrimPrefix(e.ext, ".")), size)
}

func (e folderEntry) FilterValue() string { return e.name }

// FolderBrowserScreen shows a directory listing filtered to PDF/EPUB files.
type FolderBrowserScreen struct {
	list    list.Model
	dir     string
	width   int
	height  int
	err     error
}

// NewFolderBrowserScreen creates a folder browser starting at the given directory.
func NewFolderBrowserScreen(dir string) FolderBrowserScreen {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(styles.Purple).
		BorderLeftForeground(styles.Purple)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(styles.FgSecondary).
		BorderLeftForeground(styles.Purple)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = dir
	l.Styles.Title = styles.TitleStyle.
		Background(styles.PurpleDim).
		Foreground(lipgloss.Color("#ffffff")).
		Padding(0, 2)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	s := FolderBrowserScreen{
		list: l,
		dir:  dir,
	}
	return s
}

// SetSize sets the screen dimensions.
func (s *FolderBrowserScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.list.SetSize(width, height-2)
}

func (s *FolderBrowserScreen) LoadDir() tea.Cmd {
	dir := s.dir
	return func() tea.Msg {
		entries, err := scanDir(dir)
		if err != nil {
			return folderErrMsg{err}
		}
		return folderLoadedMsg{dir: dir, entries: entries}
	}
}

type folderLoadedMsg struct {
	dir     string
	entries []folderEntry
}
type folderErrMsg struct{ err error }

func (s FolderBrowserScreen) Init() tea.Cmd {
	return s.LoadDir()
}

func (s FolderBrowserScreen) Update(msg tea.Msg) (FolderBrowserScreen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetSize(msg.Width, msg.Height-2)

	case folderLoadedMsg:
		s.dir = msg.dir
		s.list.Title = msg.dir
		items := make([]list.Item, len(msg.entries))
		for i, e := range msg.entries {
			items[i] = e
		}
		s.list.SetItems(items)
		s.list.ResetSelected()
		s.err = nil
		return s, nil

	case folderErrMsg:
		s.err = msg.err
		return s, nil

	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if item, ok := s.list.SelectedItem().(folderEntry); ok {
				if item.isDir {
					s.dir = item.path
					return s, s.LoadDir()
				}
				// It's a document — open it (add to library + open)
				path := item.path
				return s, func() tea.Msg { return FileSelectedMsg{Path: path} }
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("backspace", "-"))):
			parent := filepath.Dir(s.dir)
			if parent != s.dir {
				s.dir = parent
				return s, s.LoadDir()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("~"))):
			home, err := os.UserHomeDir()
			if err == nil {
				s.dir = home
				return s, s.LoadDir()
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "esc"))):
			return s, func() tea.Msg { return CloseFolderBrowserMsg{} }
		}
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s FolderBrowserScreen) View() string {
	help := components.HelpBar([]components.HelpBinding{
		{Key: "enter", Desc: "open"},
		{Key: "bksp", Desc: "parent dir"},
		{Key: "~", Desc: "home"},
		{Key: "/", Desc: "filter"},
		{Key: "q", Desc: "back"},
	})

	if s.err != nil {
		errStr := styles.ErrorStyle.Render("Error: " + s.err.Error())
		return lipgloss.JoinVertical(lipgloss.Left, s.list.View(), errStr, help)
	}

	return lipgloss.JoinVertical(lipgloss.Left, s.list.View(), help)
}

// Dir returns the current directory path.
func (s FolderBrowserScreen) Dir() string {
	return s.dir
}

// scanDir reads a directory and returns sorted entries (dirs first, then documents).
func scanDir(dir string) ([]folderEntry, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read directory: %w", err)
	}

	var dirs []folderEntry
	var files []folderEntry

	for _, e := range dirEntries {
		name := e.Name()
		// Skip hidden files/dirs
		if strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(dir, name)

		if e.IsDir() {
			// Only show directories that contain documents (or subdirs)
			if hasDocs, _ := dirContainsDocuments(fullPath); hasDocs {
				dirs = append(dirs, folderEntry{
					name:  name,
					path:  fullPath,
					isDir: true,
				})
			}
			continue
		}

		ext := strings.ToLower(filepath.Ext(name))
		if ext == ".pdf" || ext == ".epub" {
			info, _ := e.Info()
			var size int64
			if info != nil {
				size = info.Size()
			}
			files = append(files, folderEntry{
				name:  name,
				path:  fullPath,
				isDir: false,
				size:  size,
				ext:   ext,
			})
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].name) < strings.ToLower(dirs[j].name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].name) < strings.ToLower(files[j].name)
	})

	// Parent directory entry (unless we're at root)
	var result []folderEntry
	parent := filepath.Dir(dir)
	if parent != dir {
		result = append(result, folderEntry{
			name:  "..",
			path:  parent,
			isDir: true,
		})
	}

	result = append(result, dirs...)
	result = append(result, files...)
	return result, nil
}

// dirContainsDocuments checks if a directory (or its children, 1 level deep)
// contains any PDF or EPUB files.
func dirContainsDocuments(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".pdf" || ext == ".epub" {
			return true, nil
		}
	}
	// Check one level of subdirectories
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		subEntries, err := os.ReadDir(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		for _, se := range subEntries {
			if se.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(se.Name()))
			if ext == ".pdf" || ext == ".epub" {
				return true, nil
			}
		}
	}
	return false, nil
}

func formatSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.0f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

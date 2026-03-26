# AGENTS.md - Coding Agent Instructions for PulseReader

## Project Overview

PulseReader is a terminal-based document reader and RSS feed aggregator written in Go. It supports PDF and EPUB documents, provides a TUI built with Bubble Tea, and stores data in SQLite.

## Build/Lint/Test Commands

```bash
# Build the binary
make build
# Output: bin/pulsereader

# Run the application
make run
make run ARGS="/path/to/file.pdf"

# Install to GOPATH/bin
make install

# Clean build artifacts
make clean

# Cross-compile for all platforms
make build-all

# Run linter (if golangci-lint is installed)
golangci-lint run ./...

# Run go vet
go vet ./...

# Format code
gofmt -w .

# Run all tests (when tests exist)
go test ./...

# Run a single test
go test -run TestFunctionName ./path/to/package
go test -run TestFunctionName ./internal/storage

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

## Project Structure

```
PulseReader/
├── cmd/pulsereader/     # Main application entry point
│   └── main.go
├── internal/
│   ├── app/             # Top-level Bubble Tea model and routing
│   ├── config/          # Configuration and paths
│   ├── document/        # Document interfaces (PDF, EPUB)
│   ├── feed/            # RSS/Atom feed fetching
│   ├── storage/         # SQLite database layer
│   └── tui/
│       ├── components/  # Reusable UI components
│       ├── screens/     # Individual screen models
│       └── styles/      # Theme and styling
├── go.mod
├── go.sum
└── Makefile
```

## Code Style Guidelines

### Imports

```go
// Standard library first
import (
    "context"
    "database/sql"
    "fmt"
    "os"

    // Third-party packages second (blank line separator)
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"

    // Local packages last
    "github.com/jasonkneen/pulsereader/internal/storage"
)
```

- Group imports by category: standard library, third-party, local
- Use blank lines between import groups
- Use import aliases for Bubble Tea: `tea "github.com/charmbracelet/bubbletea"`
- Use blank identifier imports for drivers: `_ "modernc.org/sqlite"`

### Naming Conventions

- **Packages**: lowercase, single word preferred (`storage`, `feed`, `document`, `tui`)
- **Types**: PascalCase for exported, camelCase for unexported
- **Interfaces**: typically single-method interfaces with `-er` suffix (`Document`, `Reader`)
- **Structs**: PascalCase (`PDFDocument`, `LibraryScreen`)
- **Constants**: PascalCase or camelCase based on export status
- **Private fields**: lowercase first letter (`db`, `conn`, `pages`)

```go
type DB struct {
    conn *sql.DB  // private field
}

type Document interface {
    Metadata() Metadata
    TotalPages() int
    PageContent(page int) (string, error)
    Close() error
}
```

### Formatting

- Use `gofmt` (or `go fmt`) for all formatting
- Tabs for indentation (Go standard)
- No trailing whitespace
- Max line length: follow Go conventions (typically no strict limit, but be reasonable)

### Error Handling

```go
// Wrap errors with context using fmt.Errorf and %w
func Open(path string) (*DB, error) {
    conn, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, fmt.Errorf("open database: %w", err)
    }
    // ...
}

// Return early on errors
func (d *PDFDocument) PageContent(page int) (string, error) {
    if page < 1 || page > d.pages {
        return "", fmt.Errorf("page %d out of range (1-%d)", page, d.pages)
    }
    // ...
}

// For Bubble Tea commands, return error messages
type errMsg struct{ err error }

func (s *LibraryScreen) LoadDocuments() tea.Cmd {
    return func() tea.Msg {
        docs, err := s.db.ListDocuments()
        if err != nil {
            return errMsg{err}
        }
        return docsLoadedMsg{docs}
    }
}

// Discard errors explicitly with _ when intentional
_ = m.db.SetSetting("last_browse_dir", m.lastBrowseDir)
```

### Types and Interfaces

```go
// Prefer small, focused interfaces
type Document interface {
    Metadata() Metadata
    TotalPages() int
    PageContent(page int) (string, error)
    Close() error
}

// Use structs for data models with clear field types
type Article struct {
    ID          int64
    FeedID      int64
    Title       string
    PublishedAt *time.Time  // Use pointers for nullable fields
    Read        bool
}

// Use type aliases for message types in Bubble Tea
type (
    OpenDocumentMsg   struct{ Path string }
    RemoveDocumentMsg struct{ ID int64 }
    SwitchToFeedsMsg  struct{}
)
```

### Bubble Tea Patterns

```go
// Screen models follow this pattern
type LibraryScreen struct {
    list   list.Model
    db     *storage.DB
    width  int
    height int
    err    error
}

// Constructor functions start with New
func NewLibraryScreen(db *storage.DB) LibraryScreen { ... }

// Implement tea.Model interface
func (s LibraryScreen) Init() tea.Cmd { ... }
func (s LibraryScreen) Update(msg tea.Msg) (LibraryScreen, tea.Cmd) { ... }
func (s LibraryScreen) View() string { ... }

// Message types for communication between screens
type OpenDocumentMsg struct{ Path string }

// Commands that perform async operations
func (s *LibraryScreen) LoadDocuments() tea.Cmd {
    return func() tea.Msg {
        docs, err := s.db.ListDocuments()
        if err != nil {
            return errMsg{err}
        }
        return docsLoadedMsg{docs}
    }
}
```

### Database Patterns

```go
// Use defer rows.Close() for query cleanup
rows, err := db.conn.Query(query, args...)
if err != nil {
    return nil, err
}
defer rows.Close()

// Use sql.NullTime for nullable time fields
var publishedAt sql.NullTime
err := rows.Scan(&a.ID, &publishedAt, ...)
if publishedAt.Valid {
    a.PublishedAt = &publishedAt.Time
}

// Use ON CONFLICT for upserts
INSERT INTO feeds (title, url) VALUES (?, ?)
ON CONFLICT(url) DO UPDATE SET title=excluded.title
```

### Styling (lipgloss)

```go
// Define colors as lipgloss.Color
var (
    Purple = lipgloss.Color("#7c3aed")
    FgPrimary = lipgloss.Color("#e0e0e0")
)

// Create reusable styles
var TitleStyle = lipgloss.NewStyle().
    Foreground(Purple).
    Bold(true).
    Padding(0, 1)

// Compose layouts with lipgloss.JoinVertical/JoinHorizontal
return lipgloss.JoinVertical(lipgloss.Left, titleBar, content, helpBar)
```

## Key Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - Pre-built TUI components
- `github.com/charmbracelet/lipgloss` - Styling/layout
- `github.com/charmbracelet/glamour` - Markdown rendering
- `modernc.org/sqlite` - Pure Go SQLite driver
- `github.com/ledongthuc/pdf` - PDF parsing
- `github.com/kapmahc/epub` - EPUB parsing
- `github.com/mmcdole/gofeed` - RSS/Atom parsing
- `github.com/JohannesKaufmann/html-to-markdown/v2` - HTML to Markdown

## Important Notes

- This project uses pure Go SQLite (`modernc.org/sqlite`), no CGO required
- Documents use 1-indexed page numbers
- The TUI follows Bubble Tea's Elm-style architecture
- Screen routing is handled in `internal/app/app.go`
- Database migrations run automatically on startup
- User data is stored in platform-appropriate directories (XDG on Linux/macOS)

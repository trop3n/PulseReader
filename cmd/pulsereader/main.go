package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jasonkneen/pulsereader/internal/app"
	"github.com/jasonkneen/pulsereader/internal/config"
	"github.com/jasonkneen/pulsereader/internal/document"
	"github.com/jasonkneen/pulsereader/internal/storage"
	"github.com/jasonkneen/pulsereader/internal/tui/screens"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	db, err := storage.Open(cfg.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	model := app.New(db)

	// If a file path is provided as an argument, open it directly
	if len(os.Args) > 1 {
		path := os.Args[1]
		// Validate the file exists and is a supported format
		if _, err := os.Stat(path); err != nil {
			fmt.Fprintf(os.Stderr, "File not found: %s\n", path)
			os.Exit(1)
		}
		doc, err := document.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open file: %v\n", err)
			os.Exit(1)
		}
		doc.Close()

		// Queue opening the document after the TUI starts
		p := tea.NewProgram(model, tea.WithAltScreen())
		go func() {
			p.Send(screens.OpenDocumentMsg{Path: path})
		}()
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

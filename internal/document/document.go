package document

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Metadata holds common document information.
type Metadata struct {
	Title      string
	Author     string
	Format     string
	TotalPages int
	FileSize   int64
}

// Document is the interface for readable documents.
type Document interface {
	Metadata() Metadata
	TotalPages() int
	PageContent(page int) (string, error)
	Close() error
}

// Open opens a PDF or EPUB file based on its extension.
func Open(path string) (Document, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pdf":
		return OpenPDF(path)
	case ".epub":
		return OpenEPUB(path)
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}
}

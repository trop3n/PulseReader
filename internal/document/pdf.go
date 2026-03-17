package document

import (
	"fmt"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

// PDFDocument handles PDF text extraction.
type PDFDocument struct {
	path     string
	title    string
	pages    int
	fileSize int64
	reader   *pdf.Reader
	file     *os.File
}

// OpenPDF opens a PDF file for reading.
func OpenPDF(path string) (*PDFDocument, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parse pdf: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("stat pdf: %w", err)
	}

	title := strings.TrimSuffix(info.Name(), ".pdf")

	// Attempt to get title from PDF metadata
	trailer := r.Trailer()
	infoVal := trailer.Key("Info")
	if infoVal.Kind() != pdf.Null {
		t := infoVal.Key("Title")
		if t.Kind() != pdf.Null {
			s := t.String()
			s = strings.Trim(s, "()")
			if s != "" {
				title = s
			}
		}
	}

	return &PDFDocument{
		path:     path,
		title:    title,
		pages:    r.NumPage(),
		fileSize: info.Size(),
		reader:   r,
		file:     f,
	}, nil
}

// Metadata returns document metadata.
func (d *PDFDocument) Metadata() Metadata {
	return Metadata{
		Title:      d.title,
		Author:     "",
		Format:     "pdf",
		TotalPages: d.pages,
		FileSize:   d.fileSize,
	}
}

// TotalPages returns the number of pages.
func (d *PDFDocument) TotalPages() int {
	return d.pages
}

// PageContent extracts text content from a specific page (1-indexed).
func (d *PDFDocument) PageContent(page int) (string, error) {
	if page < 1 || page > d.pages {
		return "", fmt.Errorf("page %d out of range (1-%d)", page, d.pages)
	}

	p := d.reader.Page(page)
	if p.V.IsNull() {
		return "", fmt.Errorf("page %d is empty", page)
	}

	text, err := p.GetPlainText(nil)
	if err != nil {
		return "", fmt.Errorf("extract text from page %d: %w", page, err)
	}

	return cleanText(text), nil
}

// Close releases resources.
func (d *PDFDocument) Close() error {
	return d.file.Close()
}

// cleanText normalizes whitespace in extracted PDF text.
func cleanText(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	blankCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			blankCount++
			if blankCount <= 1 {
				result = append(result, "")
			}
		} else {
			blankCount = 0
			result = append(result, trimmed)
		}
	}
	return strings.TrimSpace(strings.Join(result, "\n"))
}

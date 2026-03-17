package document

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kapmahc/epub"
	md "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// EPUBDocument handles EPUB parsing and chapter extraction.
type EPUBDocument struct {
	path     string
	book     *epub.Book
	fileSize int64
	chapters []chapter
}

type chapter struct {
	title string
	href  string
}

// OpenEPUB opens an EPUB file for reading.
func OpenEPUB(path string) (*EPUBDocument, error) {
	book, err := epub.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open epub: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		book.Close()
		return nil, fmt.Errorf("stat epub: %w", err)
	}

	doc := &EPUBDocument{
		path:     path,
		book:     book,
		fileSize: info.Size(),
	}
	doc.parseChapters()

	return doc, nil
}

func (d *EPUBDocument) parseChapters() {
	// Use the spine order (reading order) from the EPUB
	for _, item := range d.book.Opf.Spine.Items {
		itemID := item.IDref
		// Find the manifest item
		for _, mi := range d.book.Opf.Manifest {
			if mi.ID == itemID {
				title := mi.ID // fallback title
				// Try to find a better title from the TOC
				for _, navPoint := range d.book.Ncx.Points {
					if strings.Contains(navPoint.Content.Src, mi.Href) {
						title = navPoint.Text
						break
					}
				}
				d.chapters = append(d.chapters, chapter{
					title: title,
					href:  mi.Href,
				})
				break
			}
		}
	}

	// Fallback: if no spine items, use manifest items that are XHTML
	if len(d.chapters) == 0 {
		for _, item := range d.book.Opf.Manifest {
			if strings.Contains(item.MediaType, "html") {
				d.chapters = append(d.chapters, chapter{
					title: item.ID,
					href:  item.Href,
				})
			}
		}
	}
}

// Metadata returns document metadata.
func (d *EPUBDocument) Metadata() Metadata {
	title := d.path
	if len(d.book.Opf.Metadata.Title) > 0 {
		title = d.book.Opf.Metadata.Title[0]
	}
	author := ""
	if len(d.book.Opf.Metadata.Creator) > 0 {
		author = d.book.Opf.Metadata.Creator[0].Data
	}
	return Metadata{
		Title:      title,
		Author:     author,
		Format:     "epub",
		TotalPages: len(d.chapters),
		FileSize:   d.fileSize,
	}
}

// TotalPages returns the number of chapters.
func (d *EPUBDocument) TotalPages() int {
	return len(d.chapters)
}

// ChapterTitles returns the title of each chapter.
func (d *EPUBDocument) ChapterTitles() []string {
	titles := make([]string, len(d.chapters))
	for i, ch := range d.chapters {
		titles[i] = ch.title
	}
	return titles
}

// PageContent returns the rendered markdown content of a chapter (1-indexed).
func (d *EPUBDocument) PageContent(page int) (string, error) {
	if page < 1 || page > len(d.chapters) {
		return "", fmt.Errorf("chapter %d out of range (1-%d)", page, len(d.chapters))
	}

	ch := d.chapters[page-1]

	// Open the chapter file from the EPUB archive
	f, err := d.book.Open(ch.href)
	if err != nil {
		return "", fmt.Errorf("open chapter %q: %w", ch.href, err)
	}
	defer f.Close()

	// Read HTML content
	htmlBytes, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("read chapter %q: %w", ch.href, err)
	}

	html := string(htmlBytes)

	// Convert HTML to Markdown
	markdown, err := md.ConvertString(html)
	if err != nil {
		// Fallback: return raw text with HTML tags stripped
		return stripTags(html), nil
	}

	return markdown, nil
}

// Close releases resources.
func (d *EPUBDocument) Close() error {
	d.book.Close()
	return nil
}

// stripTags is a simple HTML tag stripper for fallback.
func stripTags(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	return result.String()
}

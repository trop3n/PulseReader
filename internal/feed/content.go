package feed

import (
	"fmt"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/jasonkneen/pulsereader/internal/storage"
)

// RenderArticle converts an article's HTML content to terminal-friendly markdown.
// Returns a formatted string with title, metadata, and body.
func RenderArticle(article *storage.Article) string {
	var b strings.Builder

	// Title
	b.WriteString("# ")
	b.WriteString(article.Title)
	b.WriteString("\n\n")

	// Metadata line
	var meta []string
	if article.Author != "" {
		meta = append(meta, "By "+article.Author)
	}
	if article.FeedTitle != "" {
		meta = append(meta, article.FeedTitle)
	}
	if article.PublishedAt != nil {
		meta = append(meta, article.PublishedAt.Format("Jan 2, 2006"))
	}
	if len(meta) > 0 {
		b.WriteString("*")
		b.WriteString(strings.Join(meta, " · "))
		b.WriteString("*\n\n---\n\n")
	}

	// Body content — prefer full content, fall back to summary
	html := article.Content
	if html == "" {
		html = article.Summary
	}

	if html != "" {
		markdown, err := md.ConvertString(html)
		if err != nil {
			// Fallback: basic tag stripping
			b.WriteString(stripHTML(html))
		} else {
			b.WriteString(markdown)
		}
	} else {
		b.WriteString("*(no content available)*")
	}

	// Link to original
	if article.URL != "" {
		b.WriteString(fmt.Sprintf("\n\n---\n\n[Read original](%s)", article.URL))
	}

	return b.String()
}

func stripHTML(s string) string {
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

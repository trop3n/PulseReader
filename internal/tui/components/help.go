package components

import (
	"strings"

	"github.com/jasonkneen/pulsereader/internal/tui/styles"
)

// HelpBinding is a key-description pair.
type HelpBinding struct {
	Key  string
	Desc string
}

// HelpBar renders a compact help line from bindings.
func HelpBar(bindings []HelpBinding) string {
	var parts []string
	for _, b := range bindings {
		key := styles.AccentStyle.Render(b.Key)
		parts = append(parts, key+" "+b.Desc)
	}
	return styles.HelpStyle.Render(strings.Join(parts, "  "))
}

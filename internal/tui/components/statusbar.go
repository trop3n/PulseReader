package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/jasonkneen/pulsereader/internal/tui/styles"
)

// StatusBar renders a bottom status bar.
func StatusBar(width int, left, center, right string) string {
	style := styles.StatusBarStyle.Width(width)

	leftStr := lipgloss.NewStyle().Foreground(styles.Purple).Bold(true).Render(left)
	centerStr := lipgloss.NewStyle().Foreground(styles.FgSecondary).Render(center)
	rightStr := lipgloss.NewStyle().Foreground(styles.FgMuted).Render(right)

	// Calculate spacing
	leftLen := lipgloss.Width(leftStr)
	centerLen := lipgloss.Width(centerStr)
	rightLen := lipgloss.Width(rightStr)

	gap1 := (width - leftLen - centerLen - rightLen) / 2
	gap2 := width - leftLen - centerLen - rightLen - gap1

	if gap1 < 1 {
		gap1 = 1
	}
	if gap2 < 1 {
		gap2 = 1
	}

	content := fmt.Sprintf("%s%*s%s%*s%s", leftStr, gap1, "", centerStr, gap2, "", rightStr)
	return style.Render(content)
}

package styles

import "github.com/charmbracelet/lipgloss"

// Colors — dark theme with purple accent (matching original PulseReader).
var (
	Purple     = lipgloss.Color("#7c3aed")
	PurpleDim  = lipgloss.Color("#5b21b6")
	BgDark     = lipgloss.Color("#1e1e1e")
	BgPanel    = lipgloss.Color("#252525")
	BgSelected = lipgloss.Color("#2d2d2d")
	FgPrimary  = lipgloss.Color("#e0e0e0")
	FgSecondary = lipgloss.Color("#999999")
	FgMuted    = lipgloss.Color("#666666")
	Green      = lipgloss.Color("#22c55e")
	Red        = lipgloss.Color("#ef4444")
	Yellow     = lipgloss.Color("#eab308")
)

// App-level styles.
var (
	AppStyle = lipgloss.NewStyle().
			Background(BgDark)

	TitleStyle = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(FgSecondary).
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(FgPrimary).
				Background(BgSelected).
				Bold(true).
				Padding(0, 1)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(FgSecondary).
			Padding(0, 1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(FgSecondary).
			Background(BgPanel).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(FgMuted).
			Padding(0, 1)

	AccentStyle = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true).
			Padding(0, 1)

	ReaderContentStyle = lipgloss.NewStyle().
				Foreground(FgPrimary).
				Padding(1, 4)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PurpleDim)
)

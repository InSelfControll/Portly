package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors
const (
	PrimaryColor   = "#7C3AED" // Violet
	SecondaryColor = "#10B981" // Emerald
	SuccessColor   = "#22C55E" // Green
	WarningColor   = "#F59E0B" // Amber
	ErrorColor     = "#EF4444" // Red
	InfoColor      = "#3B82F6" // Blue
	TextColor      = "#E5E7EB" // Gray 200
	MutedColor     = "#9CA3AF" // Gray 400
	DarkBg         = "#111827" // Gray 900
	PanelBg        = "#1F2937" // Gray 800
	BorderColor    = "#374151" // Gray 700
)

// Styles
var (
	// Title style
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(PrimaryColor)).
		Padding(0, 1).
		MarginBottom(1)

	// Subtitle style
	Subtitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(MutedColor)).
		MarginBottom(1)

	// Panel style
	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderColor)).
		Background(lipgloss.Color(PanelBg)).
		Padding(1, 2)

	// Active panel
	ActivePanel = Panel.Copy().
		BorderForeground(lipgloss.Color(PrimaryColor))

	// Menu item
	MenuItem = lipgloss.NewStyle().
		Padding(0, 1)

	// Active menu item
	ActiveMenuItem = MenuItem.Copy().
		Bold(true).
		Foreground(lipgloss.Color(PrimaryColor)).
		Background(lipgloss.Color(DarkBg))

	// Success message
	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color(SuccessColor))

	// Error message
	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ErrorColor))

	// Warning message
	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color(WarningColor))

	// Info message
	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color(InfoColor))

	// Help text
	Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color(MutedColor))

	// Status bar
	StatusBar = lipgloss.NewStyle().
		Background(lipgloss.Color(PrimaryColor)).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	// Table header
	TableHeader = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(PrimaryColor))

	// Table cell
	TableCell = lipgloss.NewStyle().
		Padding(0, 1)

	// Button
	Button = lipgloss.NewStyle().
		Background(lipgloss.Color(PrimaryColor)).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2).
		Margin(0, 1)

	// Button focus
	ButtonFocus = Button.Copy().
		Background(lipgloss.Color(SecondaryColor))
)

// Width returns a style with specific width
func Width(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width)
}

// Center centers content
func Center(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
}

package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

// viewAddRule renders the add rule form
func (m *Model) viewAddRule() string {
	form := m.addRuleForm

	title := styles.Title.Render("Add NAT Rule")
	subtitle := styles.Subtitle.Render("Select a product to auto-fill suggested ports")

	// Get product info for displaying suggestions
	productInfo := form.fields[0].GetProductInfo()
	suggestedPort := 0
	if len(productInfo.DefaultPorts) > 0 {
		suggestedPort = productInfo.DefaultPorts[0]
	}

	var fields []string
	for i, field := range form.fields {
		label := field.label
		if field.required {
			label += " *"
		}

		style := lipgloss.NewStyle()
		if i == form.focus {
			style = style.Foreground(lipgloss.Color(styles.PrimaryColor)).Bold(true)
		} else {
			style = style.Foreground(lipgloss.Color(styles.MutedColor))
		}

		labelStr := style.Render(label + ":")
		inputStr := field.View()

		// Show suggestion for port fields
		suggestion := ""
		if suggestedPort > 0 && field.fieldType == FieldTypeNumber {
			currentVal := field.Value()
			if currentVal == "" || currentVal == strconv.Itoa(suggestedPort) {
				suggestion = lipgloss.NewStyle().
					Foreground(lipgloss.Color(styles.InfoColor)).
					Render(fmt.Sprintf(" (suggested: %d)", suggestedPort))
			}
		}

		fieldContent := lipgloss.JoinHorizontal(
			lipgloss.Left,
			labelStr,
			suggestion,
		)
		fieldContent = lipgloss.JoinVertical(
			lipgloss.Left,
			fieldContent,
			inputStr,
		)

		// Show product dropdown if active
		if i == form.focus && field.fieldType == FieldTypeProduct && field.ShowOptions() {
			fieldContent = renderProductDropdown(fieldContent, field, form.optionFocus)
		}

		fields = append(fields, fieldContent)
	}

	formContent := lipgloss.JoinVertical(lipgloss.Left, fields...)

	var help string
	if form.ShowingOptions() {
		help = styles.Help.Render("↑/↓: select • enter: confirm • tab: close")
	} else {
		help = styles.Help.Render("tab: next • enter: submit • esc: back • Ctrl+D: product list")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		styles.Panel.Render(formContent),
		"",
		help,
	)
}

// renderProductDropdown renders the product dropdown
func renderProductDropdown(fieldContent string, field EnhancedFormField, optionFocus int) string {
	var options []string
	for j, opt := range field.Options() {
		optStyle := lipgloss.NewStyle()
		if j == optionFocus {
			optStyle = optStyle.Background(lipgloss.Color(styles.PrimaryColor)).Foreground(lipgloss.Color("#FFFFFF"))
		} else {
			optStyle = optStyle.Foreground(lipgloss.Color(styles.TextColor))
		}
		options = append(options, optStyle.Render("  "+opt))
	}
	dropdown := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(styles.BorderColor)).
		Render(lipgloss.JoinVertical(lipgloss.Left, options...))
	return lipgloss.JoinVertical(lipgloss.Left, fieldContent, dropdown)
}

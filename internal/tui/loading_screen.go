package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/orchestrator/unified-firewall/internal/tui/styles"
)

func (m *Model) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.lastError = msg.err
		m.screen = ScreenError
		return m, nil
	case successMsg:
		m.successMsg = msg.msg
		m.screen = ScreenSuccess
		m.addRuleForm.Reset()
		return m, nil
	}
	return m, nil
}

func (m *Model) viewLoading() string {
	title := styles.Title.Render("Processing")
	spinner := "⏳ " + m.loadingMsg

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		styles.Info.Render(spinner),
	)
}

func (m *Model) viewError() string {
	title := styles.Title.Render("Error")
	errMsg := styles.Error.Render(m.lastError.Error())
	help := styles.Help.Render("any key: back")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		styles.Panel.Render(errMsg),
		"",
		help,
	)
}

func (m *Model) viewSuccess() string {
	title := styles.Title.Render("Success")
	success := styles.Success.Render(m.successMsg)
	help := styles.Help.Render("any key: menu")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		styles.Panel.Render(success),
		"",
		help,
	)
}

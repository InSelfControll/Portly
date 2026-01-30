package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/orchestrator/unified-firewall/internal/platform"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// errMsg represents an error message
type errMsg struct {
	err error
}

// successMsg represents a success message
type successMsg struct {
	msg string
}

// submitAddRule submits the form
func (m *Model) submitAddRule() (tea.Model, tea.Cmd) {
	if !platform.IsRoot() {
		m.lastError = fmt.Errorf("root privileges required")
		m.screen = ScreenError
		m.addRuleForm.Reset()
		return m, nil
	}

	if m.provider == nil {
		m.lastError = fmt.Errorf("no firewall provider available")
		m.screen = ScreenError
		m.addRuleForm.Reset()
		return m, nil
	}

	// 1. Get the populated rule data
	rule, _ := m.addRuleForm.GetRule()

	// 2. Generate and assign the ID FIRST
	rule.ID = uuid.New().String()[:8]

	// 3. Perform validation NOW that the ID is present
	if err := rule.Validate(); err != nil {
		m.lastError = err
		m.screen = ScreenError
		return m, nil
	}

	m.loadingMsg = "Applying NAT rule..."
	m.screen = ScreenLoading

	return m, func() tea.Msg {
		err := m.provider.ApplyNAT(m.ctx, rule)
		if err != nil {
			return errMsg{err}
		}

		if m.stateMgr != nil {
			appliedRule := models.AppliedRule{
				NATRule:   rule,
				Status:    models.StatusActive,
				AppliedAt: time.Now().UTC().Format(time.RFC3339),
			}
			m.stateMgr.AddRule(appliedRule)
		}

		return successMsg{fmt.Sprintf("Rule %s created successfully", rule.ID)}
	}
}

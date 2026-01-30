package tui

import (
	"context"
	"fmt"

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

	// 1. Get the populated rule data based on type
	var ruleID string
	var operation func(context.Context) error

	if m.addRuleForm.formType == FormTypeNAT {
		var rule models.NATRule
		rule, _ = m.addRuleForm.GetNATRule()
		rule.ID = uuid.New().String()[:8]
		ruleID = rule.ID
		if err := rule.Validate(); err != nil {
			m.lastError = err
			m.screen = ScreenError
			return m, nil
		}
		operation = func(ctx context.Context) error {
			return m.provider.ApplyNAT(ctx, rule)
		}
	} else {
		var rule models.FirewallRule
		rule, _ = m.addRuleForm.GetFirewallRule()
		rule.ID = uuid.New().String()[:8]
		ruleID = rule.ID
		if err := rule.Validate(); err != nil {
			m.lastError = err
			m.screen = ScreenError
			return m, nil
		}

		if rule.Type == models.RuleTypePortLimit {
			operation = func(ctx context.Context) error {
				return m.provider.OpenPortForIP(ctx, rule)
			}
		} else if rule.Type == models.RuleTypeTrustIP {
			operation = func(ctx context.Context) error {
				return m.provider.TrustIP(ctx, rule)
			}
		} else {
			operation = func(ctx context.Context) error {
				return m.provider.OpenPort(ctx, rule)
			}
		}
	}

	m.loadingMsg = "Applying rule..."
	m.screen = ScreenLoading

	return m, func() tea.Msg {
		err := operation(m.ctx)
		if err != nil {
			return errMsg{err}
		}

		// Since we don't have a generic AddRule to stateMgr yet for other types,
		// we might need to skip stateMgr update for non-NAT or fit it in.
		// For now we just focus on applying the rule.

		return successMsg{fmt.Sprintf("Rule %s created successfully", ruleID)}
	}
}

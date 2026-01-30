package state

import (
	"fmt"

	"github.com/orchestrator/unified-firewall/pkg/models"
)

// AddRule adds a rule to the state
func (m *Manager) AddRule(rule models.AppliedRule) error {
	if m.state.FindRule(rule.ID) != nil {
		return fmt.Errorf("rule with ID %s already exists", rule.ID)
	}

	m.state.Rules = append(m.state.Rules, rule)
	return m.Save()
}

// UpdateRule updates an existing rule
func (m *Manager) UpdateRule(rule models.AppliedRule) error {
	for i, r := range m.state.Rules {
		if r.ID == rule.ID {
			m.state.Rules[i] = rule
			return m.Save()
		}
	}
	return fmt.Errorf("rule with ID %s not found", rule.ID)
}

// RemoveRule removes a rule from the state
func (m *Manager) RemoveRule(ruleID string) error {
	if !m.state.RemoveRule(ruleID) {
		return fmt.Errorf("rule with ID %s not found", ruleID)
	}
	return m.Save()
}

// GetRule returns a rule by ID
func (m *Manager) GetRule(ruleID string) *models.AppliedRule {
	return m.state.FindRule(ruleID)
}

// ListRules returns all rules
func (m *Manager) ListRules() []models.AppliedRule {
	return m.state.Rules
}

// ListRulesByProduct returns rules for a specific product
func (m *Manager) ListRulesByProduct(product string) []models.AppliedRule {
	var result []models.AppliedRule
	for _, r := range m.state.Rules {
		if r.Product == product {
			result = append(result, r)
		}
	}
	return result
}

// ListActiveRules returns only active rules
func (m *Manager) ListActiveRules() []models.AppliedRule {
	var result []models.AppliedRule
	for _, r := range m.state.Rules {
		if r.Status == models.StatusActive {
			result = append(result, r)
		}
	}
	return result
}

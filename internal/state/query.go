package state

import "github.com/orchestrator/unified-firewall/pkg/models"

// IsRuleActive checks if a rule with the same parameters is already active
func (m *Manager) IsRuleActive(externalPort int, proto models.Protocol) bool {
	for _, r := range m.state.Rules {
		if r.Status == models.StatusActive &&
			r.ExternalPort == externalPort &&
			r.Proto == proto {
			return true
		}
	}
	return false
}

// GetRuleByPort returns a rule by external port and protocol
func (m *Manager) GetRuleByPort(port int, proto models.Protocol) *models.AppliedRule {
	for _, r := range m.state.Rules {
		if r.ExternalPort == port && r.Proto == proto {
			return &r
		}
	}
	return nil
}

// Rollback marks a rule as failed and returns the rule
func (m *Manager) Rollback(ruleID string, errMsg string) (*models.AppliedRule, error) {
	for i, r := range m.state.Rules {
		if r.ID == ruleID {
			m.state.Rules[i].Status = models.StatusFailed
			m.state.Rules[i].ErrorMsg = errMsg
			if err := m.Save(); err != nil {
				return nil, err
			}
			return &m.state.Rules[i], nil
		}
	}
	return nil, nil
}

// Cleanup removes failed rules from state
func (m *Manager) Cleanup() error {
	var activeRules []models.AppliedRule
	for _, r := range m.state.Rules {
		if r.Status != models.StatusRemoved {
			activeRules = append(activeRules, r)
		}
	}
	m.state.Rules = activeRules
	return m.Save()
}

// SetProductInfo updates product information
func (m *Manager) SetProductInfo(info models.ProductInfo) error {
	m.state.Products[info.Name] = info
	return m.Save()
}

// GetProductInfo returns product information
func (m *Manager) GetProductInfo(name string) (models.ProductInfo, bool) {
	info, ok := m.state.Products[name]
	return info, ok
}

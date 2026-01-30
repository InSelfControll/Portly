package models

// State represents the current state of the orchestrator
type State struct {
	Version     string                 `yaml:"version" json:"version"`
	OS          OSInfo                 `yaml:"os" json:"os"`
	Rules       []AppliedRule          `yaml:"rules" json:"rules"`
	Products    map[string]ProductInfo `yaml:"products" json:"products"`
	LastUpdated string                 `yaml:"last_updated" json:"last_updated"`
}

// FindRule finds a rule by ID
func (s *State) FindRule(id string) *AppliedRule {
	for i := range s.Rules {
		if s.Rules[i].ID == id {
			return &s.Rules[i]
		}
	}
	return nil
}

// RemoveRule removes a rule by ID
func (s *State) RemoveRule(id string) bool {
	for i, r := range s.Rules {
		if r.ID == id {
			s.Rules = append(s.Rules[:i], s.Rules[i+1:]...)
			return true
		}
	}
	return false
}

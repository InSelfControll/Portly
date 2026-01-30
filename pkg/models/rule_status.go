package models

// RuleStatus represents the current state of a rule
type RuleStatus string

const (
	StatusActive  RuleStatus = "active"
	StatusPending RuleStatus = "pending"
	StatusFailed  RuleStatus = "failed"
	StatusRemoved RuleStatus = "removed"
)

// AppliedRule tracks a rule that has been applied to the system
type AppliedRule struct {
	NATRule   `yaml:",inline" json:",inline"`
	Status    RuleStatus `yaml:"status" json:"status"`
	AppliedAt string     `yaml:"applied_at" json:"applied_at"`
	ErrorMsg  string     `yaml:"error_msg,omitempty" json:"error_msg,omitempty"`
}

package models

import (
	"fmt"
	"net"
)

// FirewallRuleType represents the type of firewall rule
type FirewallRuleType string

const (
	RuleTypeNAT       FirewallRuleType = "nat"
	RuleTypePort      FirewallRuleType = "port"
	RuleTypePortLimit FirewallRuleType = "port_limit"
)

// FirewallRule represents a firewall rule (port open, NAT, or IP-limited)
type FirewallRule struct {
	ID          string           `yaml:"id" json:"id"`
	Type        FirewallRuleType `yaml:"type" json:"type"`
	Port        int              `yaml:"port" json:"port"`
	Protocol    Protocol         `yaml:"protocol" json:"protocol"`
	SourceIP    string           `yaml:"source_ip,omitempty" json:"source_ip,omitempty"`
	Destination string           `yaml:"destination,omitempty" json:"destination,omitempty"`
	Description string           `yaml:"description" json:"description"`
	Product     string           `yaml:"product" json:"product"`
}

// Validate checks if the firewall rule is valid
func (r *FirewallRule) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("rule ID is required")
	}
	if r.Port < 1 || r.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if r.Protocol != TCP && r.Protocol != UDP {
		return fmt.Errorf("protocol must be 'tcp' or 'udp'")
	}
	if r.SourceIP != "" && net.ParseIP(r.SourceIP) == nil {
		return fmt.Errorf("source IP '%s' is not valid", r.SourceIP)
	}
	return nil
}

// String returns a human-readable representation
func (r *FirewallRule) String() string {
	if r.Type == RuleTypePortLimit && r.SourceIP != "" {
		return fmt.Sprintf("%s: %s %d/%s from %s", r.Type, r.Product, r.Port, r.Protocol, r.SourceIP)
	}
	return fmt.Sprintf("%s: %s %d/%s", r.Type, r.Product, r.Port, r.Protocol)
}

// IsPortOpen returns true if this is a simple port opening rule
func (r *FirewallRule) IsPortOpen() bool {
	return r.Type == RuleTypePort
}

// IsIPLimited returns true if this rule limits by source IP
func (r *FirewallRule) IsIPLimited() bool {
	return r.Type == RuleTypePortLimit && r.SourceIP != ""
}

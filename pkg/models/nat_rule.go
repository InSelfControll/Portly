package models

import (
	"fmt"
	"net"
)

// NATRule represents a single NAT/port forwarding rule
type NATRule struct {
	ID           string   `yaml:"id" json:"id"`
	Product      string   `yaml:"product" json:"product"`
	ExternalPort int      `yaml:"external_port" json:"external_port"`
	InternalIP   string   `yaml:"internal_ip" json:"internal_ip"`
	InternalPort int      `yaml:"internal_port" json:"internal_port"`
	Proto        Protocol `yaml:"protocol" json:"protocol"`
	Description  string   `yaml:"description" json:"description"`
}

// Validate checks if the NAT rule has valid fields
func (r *NATRule) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("NAT rule ID is required")
	}
	if r.Product == "" {
		return fmt.Errorf("product name is required")
	}
	if r.ExternalPort < 1 || r.ExternalPort > 65535 {
		return fmt.Errorf("external port must be between 1 and 65535")
	}
	if r.InternalPort < 1 || r.InternalPort > 65535 {
		return fmt.Errorf("internal port must be between 1 and 65535")
	}
	if net.ParseIP(r.InternalIP) == nil {
		return fmt.Errorf("internal IP '%s' is not valid", r.InternalIP)
	}
	if r.Proto != TCP && r.Proto != UDP {
		return fmt.Errorf("protocol must be 'tcp' or 'udp'")
	}
	return nil
}

// String returns a human-readable representation of the rule
func (r *NATRule) String() string {
	return fmt.Sprintf("%s: %s (%d) -> %s:%d/%s",
		r.Product, r.ID, r.ExternalPort, r.InternalIP, r.InternalPort, r.Proto)
}

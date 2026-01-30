package models

// SecurityPolicy represents security configuration for a product
type SecurityPolicy struct {
	SelinuxBooleans []string `yaml:"selinux_booleans" json:"selinux_booleans"`
	AppArmorProfile string   `yaml:"apparmor_profile_path" json:"apparmor_profile_path"`
	RequiresRoot    bool     `yaml:"requires_root" json:"requires_root"`
}

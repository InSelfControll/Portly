package models

// OSInfo contains information about the operating system
type OSInfo struct {
	Family       string `json:"family"`        // "rhel", "debian", "darwin"
	Distribution string `json:"distribution"`  // "fedora", "ubuntu", "macos"
	Version      string `json:"version"`
	Codename     string `json:"codename"`
}

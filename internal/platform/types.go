package platform

// OSFamily represents the operating system family
type OSFamily string

const (
	FamilyRHEL    OSFamily = "rhel"
	FamilyDebian  OSFamily = "debian"
	FamilyDarwin  OSFamily = "darwin"
	FamilyUnknown OSFamily = "unknown"
)

// OSInfo contains detected OS information
type OSInfo struct {
	Family       OSFamily
	Distribution string
	Version      string
	Codename     string
	ID           string
	IDLike       string
}

// IsRHEL returns true if the OS is RHEL family
func (o *OSInfo) IsRHEL() bool {
	return o.Family == FamilyRHEL
}

// IsDebian returns true if the OS is Debian family
func (o *OSInfo) IsDebian() bool {
	return o.Family == FamilyDebian
}

// IsDarwin returns true if the OS is macOS
func (o *OSInfo) IsDarwin() bool {
	return o.Family == FamilyDarwin
}

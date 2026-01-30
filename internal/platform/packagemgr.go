package platform

// PackageManager returns the appropriate package manager for the OS
func (o *OSInfo) PackageManager() string {
	switch o.Family {
	case FamilyRHEL:
		return "dnf"
	case FamilyDebian:
		return "apt"
	case FamilyDarwin:
		return "brew"
	default:
		return ""
	}
}

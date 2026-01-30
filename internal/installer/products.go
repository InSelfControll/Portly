package installer

// SupportedProducts defines configurations for supported products
var SupportedProducts = map[string]ProductConfig{
	"docker": {
		Name:         "docker",
		DisplayName:  "Docker",
		PackageName:  "docker-ce",
		BrewFormula:  "docker",
		Dependencies: []string{"docker-compose"},
	},
	"podman": {
		Name:         "podman",
		DisplayName:  "Podman",
		PackageName:  "podman",
		BrewFormula:  "podman",
		Dependencies: []string{"podman-compose"},
	},
	"tailscale": {
		Name:        "tailscale",
		DisplayName: "Tailscale",
		PackageName: "tailscale",
		BrewFormula: "tailscale",
	},
	"headscale": {
		Name:        "headscale",
		DisplayName: "Headscale",
		PackageName: "headscale",
		BrewFormula: "headscale",
	},
	"twingate": {
		Name:        "twingate",
		DisplayName: "Twingate",
		PackageName: "twingate",
		BrewFormula: "twingate",
	},
}

// GetSupportedProducts returns a list of supported product names
func GetSupportedProducts() []string {
	products := make([]string, 0, len(SupportedProducts))
	for name := range SupportedProducts {
		products = append(products, name)
	}
	return products
}

// IsProductSupported checks if a product is in the supported list
func IsProductSupported(name string) bool {
	_, ok := SupportedProducts[name]
	return ok
}

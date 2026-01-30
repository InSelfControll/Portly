package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

// ProductInfo holds default configuration for a product
type ProductInfo struct {
	Name         string
	Description  string
	DefaultPorts []int
}

// ProductDatabase contains default configurations for known products
var ProductDatabase = map[string]ProductInfo{
	"podman":    {Name: "podman", Description: "Container engine", DefaultPorts: []int{8080, 8443}},
	"docker":    {Name: "docker", Description: "Container platform", DefaultPorts: []int{8080, 443}},
	"tailscale": {Name: "tailscale", Description: "VPN mesh network", DefaultPorts: []int{41641}},
	"headscale": {Name: "headscale", Description: "Self-hosted Tailscale", DefaultPorts: []int{8080}},
	"twingate":  {Name: "twingate", Description: "Zero trust network", DefaultPorts: []int{443}},
	"steam":     {Name: "steam", Description: "Gaming platform", DefaultPorts: []int{27015}},
	"minecraft": {Name: "minecraft", Description: "Minecraft server", DefaultPorts: []int{25565}},
	"nginx":     {Name: "nginx", Description: "Web server", DefaultPorts: []int{80, 443}},
	"postgres":  {Name: "postgres", Description: "PostgreSQL database", DefaultPorts: []int{5432}},
	"redis":     {Name: "redis", Description: "Redis cache", DefaultPorts: []int{6379}},
	"custom":    {Name: "custom", Description: "Custom product", DefaultPorts: []int{}},
}

// NewProductField creates a product selector field with default options
func NewProductField(label string, required bool) EnhancedFormField {
	input := textinput.New()
	input.Placeholder = "Select or type custom product"

	// Build display options from database
	var options []string
	for key, info := range ProductDatabase {
		if key == "custom" {
			options = append(options, "custom       - Type your own product name")
		} else {
			portsStr := formatPorts(info.DefaultPorts)
			option := fmt.Sprintf("%-12s - %s (ports: %s)", key, info.Description, portsStr)
			options = append(options, option)
		}
	}

	return EnhancedFormField{
		label:       label,
		input:       input,
		required:    required,
		fieldType:   FieldTypeProduct,
		options:     options,
		showOptions: false,
	}
}

// formatPorts formats a list of ports for display
func formatPorts(ports []int) string {
	if len(ports) == 0 {
		return "any"
	}
	var parts []string
	for _, p := range ports {
		parts = append(parts, strconv.Itoa(p))
	}
	return strings.Join(parts, ", ")
}

// GetProductInfo extracts product info from the selected value
func (f *EnhancedFormField) GetProductInfo() ProductInfo {
	val := f.Value()

	// Check if it's a selection from the dropdown
	for _, opt := range f.options {
		if opt == val {
			// Extract product name (first word)
			fields := strings.Fields(opt)
			if len(fields) > 0 {
				name := fields[0]
				if info, ok := ProductDatabase[name]; ok {
					return info
				}
			}
		}
	}

	// Check if it's a direct product name
	if info, ok := ProductDatabase[val]; ok {
		return info
	}

	// Custom product
	return ProductInfo{Name: val, Description: "Custom product", DefaultPorts: []int{}}
}

// GetProductName extracts clean product name
func (f *EnhancedFormField) GetProductName() string {
	return f.GetProductInfo().Name
}

// GetFirstPort returns the first suggested port for this product
func (f *EnhancedFormField) GetFirstPort() int {
	info := f.GetProductInfo()
	if len(info.DefaultPorts) > 0 {
		return info.DefaultPorts[0]
	}
	return 0
}

// ToggleOptions shows/hides product dropdown
func (f *EnhancedFormField) ToggleOptions() {
	f.showOptions = !f.showOptions
}

// SelectOption selects an option by index
func (f *EnhancedFormField) SelectOption(index int) {
	if index >= 0 && index < len(f.options) {
		f.input.SetValue(f.options[index])
		f.showOptions = false
	}
}

// Options returns available options
func (f *EnhancedFormField) Options() []string {
	return f.options
}

// ShowOptions returns if options should be shown
func (f *EnhancedFormField) ShowOptions() bool {
	return f.showOptions
}

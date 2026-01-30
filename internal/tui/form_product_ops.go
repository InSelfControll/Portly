package tui

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// CheckProductChange checks if product changed and updates ports
func (f *AddRuleForm) CheckProductChange() bool {
	currentProduct := f.fields[0].GetProductName()
	if currentProduct != f.lastProduct && currentProduct != "" {
		f.lastProduct = currentProduct
		info := f.fields[0].GetProductInfo()

		// Auto-populate ports if we have defaults and fields are empty
		if len(info.DefaultPorts) > 0 {
			// Set external port if empty
			if f.fields[1].Value() == "" {
				f.fields[1].SetValue(strconv.Itoa(info.DefaultPorts[0]))
			}
			// Set internal port if empty
			if f.fields[3].Value() == "" {
				f.fields[3].SetValue(strconv.Itoa(info.DefaultPorts[0]))
			}
		}
		return true
	}
	return false
}

// NextField moves to the next field
func (f *AddRuleForm) NextField() tea.Cmd {
	f.fields[f.focus].Blur()
	f.focus = (f.focus + 1) % len(f.fields)
	f.optionFocus = -1
	return f.fields[f.focus].Focus()
}

// IsProductField returns true if current field is product selector
func (f *AddRuleForm) IsProductField() bool {
	return f.fields[f.focus].fieldType == FieldTypeProduct
}

// ShowProductOptions shows/hides product dropdown
func (f *AddRuleForm) ShowProductOptions() {
	f.fields[f.focus].ToggleOptions()
}

// SelectProductOption selects a product from dropdown
func (f *AddRuleForm) SelectProductOption(index int) {
	f.fields[f.focus].SelectOption(index)
	f.optionFocus = -1
	// Check for auto-population after selection
	f.CheckProductChange()
}

// GetProductOptions returns product options for current field
func (f *AddRuleForm) GetProductOptions() []string {
	return f.fields[f.focus].Options()
}

// ShowingOptions returns if dropdown is visible
func (f *AddRuleForm) ShowingOptions() bool {
	return f.fields[f.focus].ShowOptions()
}

// ValidateCurrentField validates the current field
func (f *AddRuleForm) ValidateCurrentField() (string, bool) {
	field := &f.fields[f.focus]
	if field.required && field.Value() == "" {
		return field.label + " is required", false
	}
	if field.fieldType == FieldTypeNumber {
		if _, valid := field.ValidateNumber(); !valid {
			return field.label + " must be a number", false
		}
	}
	return "", true
}

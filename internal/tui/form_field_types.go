package tui

import (
	"strconv"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FieldType represents the type of form field
type FieldType int

const (
	FieldTypeText FieldType = iota
	FieldTypeNumber
	FieldTypeProduct
)

// EnhancedFormField represents a form field with validation
type EnhancedFormField struct {
	label       string
	input       textinput.Model
	required    bool
	fieldType   FieldType
	options     []string
	showOptions bool
}

// NewEnhancedFormField creates a new form field
func NewEnhancedFormField(label string, required bool, fieldType FieldType, placeholder string) EnhancedFormField {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Width = 40
	return EnhancedFormField{
		label:     label,
		input:     input,
		required:  required,
		fieldType: fieldType,
	}
}

// SetValue sets the field value
func (f *EnhancedFormField) SetValue(value string) {
	f.input.SetValue(value)
}

// Value returns the current value
func (f *EnhancedFormField) Value() string {
	return f.input.Value()
}

// Focus focuses the field
func (f *EnhancedFormField) Focus() tea.Cmd {
	return f.input.Focus()
}

// Blur removes focus
func (f *EnhancedFormField) Blur() {
	f.input.Blur()
}

// View returns the field view
func (f *EnhancedFormField) View() string {
	return f.input.View()
}

// HandleInput processes input with validation
func (f *EnhancedFormField) HandleInput(msg tea.KeyMsg) (blocked bool, cmd tea.Cmd) {
	switch msg.Type {
	case tea.KeyBackspace:
		return false, nil
	case tea.KeyRunes:
		if f.fieldType == FieldTypeNumber {
			for _, r := range msg.Runes {
				if !unicode.IsDigit(r) {
					return true, nil
				}
			}
		}
		return false, nil
	case tea.KeySpace:
		if f.fieldType == FieldTypeNumber {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}

// ValidateNumber checks if the field contains a valid number
func (f *EnhancedFormField) ValidateNumber() (int, bool) {
	if f.Value() == "" {
		return 0, !f.required
	}
	num, err := strconv.Atoi(f.Value())
	if err != nil {
		return 0, false
	}
	return num, true
}

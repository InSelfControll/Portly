package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// AddRuleForm represents the add rule form
type AddRuleForm struct {
	fields      []EnhancedFormField
	focus       int
	completed   bool
	cancelled   bool
	optionFocus int
	lastProduct string // Track last selected product to detect changes
}

// NewAddRuleForm creates a new add rule form
func NewAddRuleForm() *AddRuleForm {
	productField := NewProductField("Product", true)
	portField := NewEnhancedFormField("External Port", true, FieldTypeNumber, "8080")
	ipField := NewEnhancedFormField("Internal IP", true, FieldTypeText, "10.88.0.1")
	ipField.SetValue("127.0.0.1")
	internalPortField := NewEnhancedFormField("Internal Port", true, FieldTypeNumber, "80")
	protoField := NewEnhancedFormField("Protocol (tcp/udp)", true, FieldTypeText, "tcp")
	protoField.SetValue("tcp")
	descField := NewEnhancedFormField("Description", false, FieldTypeText, "Optional description")

	fields := []EnhancedFormField{
		productField, portField, ipField,
		internalPortField, protoField, descField,
	}

	form := &AddRuleForm{
		fields:      fields,
		focus:       0,
		optionFocus: -1,
		lastProduct: "",
	}
	fields[0].Focus()
	return form
}

// Init initializes the form
func (f *AddRuleForm) Init() tea.Cmd {
	return textinput.Blink
}

// GetRule returns the rule from form data
func (f *AddRuleForm) GetRule() (models.NATRule, error) {
	externalPort, _ := f.fields[1].ValidateNumber()
	internalPort, _ := f.fields[3].ValidateNumber()

	proto := models.TCP
	if f.fields[4].Value() == "udp" {
		proto = models.UDP
	}

	rule := models.NATRule{
		Product:      f.fields[0].GetProductName(),
		ExternalPort: externalPort,
		InternalIP:   f.fields[2].Value(),
		InternalPort: internalPort,
		Proto:        proto,
		Description:  f.fields[5].Value(),
	}
	return rule, rule.Validate()
}

// Reset clears the form
func (f *AddRuleForm) Reset() {
	f.fields[0].SetValue("")
	f.fields[1].SetValue("")
	f.fields[2].SetValue("127.0.0.1")
	f.fields[3].SetValue("")
	f.fields[4].SetValue("tcp")
	f.fields[5].SetValue("")
	f.focus = 0
	f.optionFocus = -1
	f.lastProduct = ""
	f.fields[0].Focus()
	f.completed = false
	f.cancelled = false
}

// FocusedField returns the currently focused field
func (f *AddRuleForm) FocusedField() *EnhancedFormField {
	return &f.fields[f.focus]
}

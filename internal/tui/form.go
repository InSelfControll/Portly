package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/orchestrator/unified-firewall/pkg/models"
)

// FormType represents the type of form to display
type FormType int

const (
	FormTypeNAT FormType = iota
	FormTypeOpenPort
	FormTypeOpenIPPort
	FormTypeOpenIP
)

// AddRuleForm represents the add rule form
type AddRuleForm struct {
	formType    FormType
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
	portField := NewEnhancedFormField("Port / Ext Port", true, FieldTypeNumber, "8080")
	ipField := NewEnhancedFormField("Internal IP", true, FieldTypeText, "10.88.0.1")
	ipField.SetValue("127.0.0.1")
	internalPortField := NewEnhancedFormField("Internal Port", true, FieldTypeNumber, "80")
	protoField := NewEnhancedFormField("Protocol (tcp/udp)", true, FieldTypeText, "tcp")
	protoField.SetValue("tcp")
	sourceIPField := NewEnhancedFormField("Source IP", true, FieldTypeText, "0.0.0.0")
	descField := NewEnhancedFormField("Description", false, FieldTypeText, "Optional description")

	// Order matters for indexing
	fields := []EnhancedFormField{
		productField,      // 0
		portField,         // 1: External Port (NAT) or Port (Open)
		ipField,           // 2: Internal IP (NAT only)
		internalPortField, // 3: Internal Port (NAT only)
		protoField,        // 4: Protocol
		sourceIPField,     // 5: Source IP (IP Restricted only)
		descField,         // 6: Description
	}

	form := &AddRuleForm{
		formType:    FormTypeNAT,
		fields:      fields,
		focus:       0,
		optionFocus: -1,
		lastProduct: "",
	}
	fields[0].Focus()
	return form
}

// SetType configures the form for a specific rule type
func (f *AddRuleForm) SetType(t FormType) {
	f.formType = t
	f.Reset()

	// Update field labels based on type
	if t == FormTypeNAT {
		f.fields[1].label = "External Port"
	} else {
		f.fields[1].label = "Port"
	}

	// Set focus to first visible field
	firstVisible := f.nextFocusableIndex(-1)
	if firstVisible >= 0 {
		f.fields[f.focus].Blur()
		f.focus = firstVisible
		f.fields[f.focus].Focus()
	}
}

// isFieldVisible returns true if field at index i should be shown
func (f *AddRuleForm) isFieldVisible(i int) bool {
	switch f.formType {
	case FormTypeNAT:
		// Show: Product(0), ExtPort(1), IntIP(2), IntPort(3), Proto(4), Desc(6)
		// Hide: SourceIP(5)
		return i != 5
	case FormTypeOpenPort:
		// Show: Product(0), Port(1), Proto(4), Desc(6)
		// Hide: IntIP(2), IntPort(3), SourceIP(5)
		return i == 0 || i == 1 || i == 4 || i == 6
	case FormTypeOpenIPPort:
		// Show: Product(0), Port(1), Proto(4), SourceIP(5), Desc(6)
		// Hide: IntIP(2), IntPort(3)
		return i == 0 || i == 1 || i == 4 || i == 5 || i == 6
	case FormTypeOpenIP:
		// Show: SourceIP(5), Desc(6)
		// Hide: Product(0), Port(1), IntIP(2), IntPort(3), Proto(4)
		return i == 5 || i == 6
	}
	return true
}

// nextFocusableIndex finds the next visible field index
func (f *AddRuleForm) nextFocusableIndex(current int) int {
	for i := current + 1; i < len(f.fields); i++ {
		if f.isFieldVisible(i) {
			return i
		}
	}
	return -1 // No more fields
}

// prevFocusableIndex finds the previous visible field index
func (f *AddRuleForm) prevFocusableIndex(current int) int {
	for i := current - 1; i >= 0; i-- {
		if f.isFieldVisible(i) {
			return i
		}
	}
	return -1 // No previous fields
}

// Init initializes the form
func (f *AddRuleForm) Init() tea.Cmd {
	return textinput.Blink
}

// GetNATRule returns the NAT rule from form data
func (f *AddRuleForm) GetNATRule() (models.NATRule, error) {
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
		Description:  f.fields[6].Value(), // Index changed to 6
	}
	// Return rule without calling .Validate() here
	return rule, nil
}

// GetFirewallRule returns the firewall rule from form data
func (f *AddRuleForm) GetFirewallRule() (models.FirewallRule, error) {
	port, _ := f.fields[1].ValidateNumber()

	proto := models.TCP
	if f.fields[4].Value() == "udp" {
		proto = models.UDP
	}

	ruleType := models.RuleTypePort
	if f.formType == FormTypeOpenIPPort {
		ruleType = models.RuleTypePortLimit
	} else if f.formType == FormTypeOpenIP {
		ruleType = models.RuleTypeTrustIP
	}

	rule := models.FirewallRule{
		Product:     f.fields[0].GetProductName(),
		Type:        ruleType,
		Port:        port,
		Protocol:    proto,
		SourceIP:    f.fields[5].Value(),
		Description: f.fields[6].Value(),
	}
	return rule, nil
}

// Reset clears the form
func (f *AddRuleForm) Reset() {
	f.fields[0].SetValue("")
	f.fields[1].SetValue("")
	f.fields[2].SetValue("127.0.0.1")
	f.fields[3].SetValue("")
	f.fields[4].SetValue("tcp")
	f.fields[5].SetValue("")
	f.fields[6].SetValue("")
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
